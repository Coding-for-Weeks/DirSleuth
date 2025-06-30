package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Coding-for-Weeks/dirsleuth/internal/worker"
)

func main() {
	var domain, wordlist string
	var threads int
	var useHTTPS bool

	flag.StringVar(&domain, "d", "", "Target domain")
	flag.StringVar(&wordlist, "w", "", "Wordlist file")
	flag.IntVar(&threads, "t", 10, "Number of threads")
	flag.BoolVar(&useHTTPS, "https", false, "Use HTTPS")
	flag.Parse()

	if domain == "" {
		log.Fatal("Error: Target domain is required.")
	}
	if wordlist == "" {
		log.Fatal("Error: Wordlist file is required.")
	}

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatalf("Error opening wordlist file: %s\n", err)
	}
	defer file.Close()

	urls := make(chan string, threads)
	results := make(chan string, threads)
	var wg sync.WaitGroup
	quit := make(chan struct{})
	var closeOnce sync.Once

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Start workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker.Worker(client, urls, &wg, results, quit)
	}

	// Read the wordlist and enqueue URLs
	go func() {
		defer close(urls)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			scheme := "http"
			if useHTTPS {
				scheme = "https"
			}
			select {
			case urls <- fmt.Sprintf("%s://%s/%s", scheme, domain, scanner.Text()):
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

	// Print results
	for result := range results {
		fmt.Println(result)
	}

	// Signal quit to goroutines in case of early exit
	closeOnce.Do(func() { close(quit) })
}
