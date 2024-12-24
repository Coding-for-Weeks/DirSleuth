package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "valid") {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	urls := make(chan string, 2)
	results := make(chan string, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go worker(client, urls, &wg, results)

	// Send test URLs to the worker
	urls <- server.URL + "/valid"
	urls <- server.URL + "/invalid"
	close(urls)

	// Wait for the worker to finish
	wg.Wait()
	close(results)

	// Check results
	expected := server.URL + "/valid"
	select {
	case result := <-results:
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	default:
		t.Error("Expected result, but got none")
	}
}
