package worker

import (
	"log"
	"net/http"
	"sync"
)

func Worker(client *http.Client, urls <-chan string, wg *sync.WaitGroup, results chan<- string, quit <-chan struct{}, verbose bool) {
	defer wg.Done()
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				return
			}
			resp, err := client.Get(url)
			if err != nil {
				if verbose {
					log.Printf("❌ [ERROR] %s → %s\n", url, err)
				}
				continue
			}
			if verbose {
				log.Printf("[%d] %s\n", resp.StatusCode, url)
			}
			if resp.StatusCode == 200 {
				results <- url
			}
			resp.Body.Close()
		case <-quit:
			return
		}
	}
}
