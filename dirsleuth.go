package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func worker(urls <-chan string, wg *sync.WaitGroup, results chan<- string) {
	for url := range urls {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			results <- url
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	wg.Done()
}

func main() {
	var domain, wordlist string
	var threads int

	flag.StringVar(&domain, "u", "", "Target domain")
	flag.StringVar(&wordlist, "w", "", "Wordlist file")
	flag.IntVar(&threads, "t", 10, "Number of threads")
	flag.Parse()

	if domain == "" {
		fmt.Println("Error: Target domain is required.")
		os.Exit(1)
	}
	if wordlist == "" {
		fmt.Println("Error: Wordlist file is required.")
		os.Exit(1)
	}

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Printf("Error opening wordlist file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	urls := make(chan string, threads)
	results := make(chan string)
	var wg sync.WaitGroup
	quit := make(chan struct{})

	// Start workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(urls, &wg, results)
	}

	// Read the wordlist and enqueue URLs
	go func() {
		defer close(urls)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case urls <- fmt.Sprintf("http://%s/%s", domain, scanner.Text()):
			case <-quit:
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading wordlist: %s\n", err)
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
