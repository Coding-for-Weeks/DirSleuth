package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
	"strings"
	"encoding/json"

	"github.com/Coding-for-Weeks/dirsleuth/internal/worker"
)

func isValidDomain(domain string) bool {
	// Basic domain validation
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9.-]+$`, domain)
	return matched
}

func main() {
	var domain, wordlist, userAgent, outputFormat string
	var threads, timeout int
	var useHTTPS, verbose bool
	var statusCodes string

	flag.StringVar(&domain, "d", "", "Target domain")
	flag.StringVar(&wordlist, "w", "", "Wordlist file")
	flag.IntVar(&threads, "t", 10, "Number of threads")
	flag.BoolVar(&useHTTPS, "https", false, "Use HTTPS")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.IntVar(&timeout, "timeout", 30, "HTTP request timeout in seconds")
	flag.StringVar(&userAgent, "user-agent", "DirSleuth/1.0", "Custom User-Agent header")
	flag.StringVar(&statusCodes, "status", "200", "Comma-separated HTTP status codes to report (e.g., 200,301,403)")
	flag.StringVar(&outputFormat, "output", "text", "Output format: text or json")
	flag.Parse()

	if domain == "" || !isValidDomain(domain) {
		log.Fatal("Error: Valid target domain is required.")
	}
	if wordlist == "" {
		log.Fatal("Error: Wordlist file is required.")
	}

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatalf("Error opening wordlist file: %s\n", err)
	}
	defer file.Close()

	// Parse status codes
	codeMap := make(map[int]bool)
	for _, codeStr := range strings.Split(statusCodes, ",") {
		codeStr = strings.TrimSpace(codeStr)
		if codeStr == "" { continue }
		var code int
		fmt.Sscanf(codeStr, "%d", &code)
		if code > 0 { codeMap[code] = true }
	}

	urls := make(chan string, threads)
	results := make(chan worker.Result, threads)
	var wg sync.WaitGroup
	quit := make(chan struct{})
	var closeOnce sync.Once

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt. Shutting down...")
		closeOnce.Do(func() { close(quit) })
	}()

	// Start workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker.Worker(client, urls, &wg, results, quit, verbose, userAgent, codeMap)
	}

	// Read the wordlist and enqueue URLs
	go func() {
		defer close(urls)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			path := strings.TrimSpace(scanner.Text())
			// Skip empty lines, full URLs, localhost, and dangerous patterns
			if path == "" || strings.Contains(path, "http://") || strings.Contains(path, "https://") || strings.Contains(path, "localhost") || strings.Contains(path, "..") {
				continue
			}
			scheme := "http"
			if useHTTPS { scheme = "https" }
			select {
			case urls <- fmt.Sprintf("%s://%s/%s", scheme, domain, path):
			case <-quit:
				return
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading wordlist: %s\n", err)
			closeOnce.Do(func() { close(quit) })
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	var output []worker.Result
	for result := range results {
		if outputFormat == "json" {
			output = append(output, result)
		} else {
			fmt.Printf("[%d] %s\n", result.StatusCode, result.URL)
		}
	}

	if outputFormat == "json" {
		json.NewEncoder(os.Stdout).Encode(output)
	}

	closeOnce.Do(func() { close(quit) })
}
