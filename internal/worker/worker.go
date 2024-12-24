package worker

import (
	"log"
	"net/http"
	"sync"
)

func Worker(client *http.Client, urls <-chan string, wg *sync.WaitGroup, results chan<- string, quit <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				return
			}
			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Error fetching URL %s: %s\n", url, err)
				continue
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
