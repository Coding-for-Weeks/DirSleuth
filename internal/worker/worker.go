package worker

import (
	"log"
	"net/http"
	"sync"
)

type Result struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

func Worker(client *http.Client, urls <-chan string, wg *sync.WaitGroup, results chan<- Result, quit <-chan struct{}, verbose bool, userAgent string, codeMap map[int]bool) {
	defer wg.Done()
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				return
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				if verbose {
					log.Printf("❌ [ERROR] %s → %s\n", url, err)
				}
				continue
			}
			req.Header.Set("User-Agent", userAgent)
			resp, err := client.Do(req)
			if err != nil {
				if verbose {
					log.Printf("❌ [ERROR] %s → %s\n", url, err)
				}
				continue
			}
			defer resp.Body.Close()
			if verbose {
				log.Printf("[%d] %s\n", resp.StatusCode, url)
			}
			if codeMap[resp.StatusCode] {
				results <- Result{URL: url, StatusCode: resp.StatusCode}
			}
		case <-quit:
			return
		}
	}
}
