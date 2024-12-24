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
)

func worker(client *http.Client, urls <-chan string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	for url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Error fetching URL %s: %s\n", url, err)
			continue
		}
		if resp.StatusCode == 200 {
			results <- url
		}
		resp.Body.Close()
	}
}

func main() {
	var domain, wordlist string
	var threads int
	var useHTTPS bool

	flag.StringVar(&domain, "u", "", "Target domain")
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Start workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(client, urls, &wg, results)
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
			close(quit)
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
	close(quit)
}
