package worker_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Coding-for-Weeks/dirsleuth/internal/worker"
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
	quit := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go worker.Worker(client, urls, &wg, results, quit)

	// Send test URLs to the worker
	urls <- server.URL + "/valid"
	urls <- server.URL + "/invalid"
	close(urls)

	// Wait for the worker to finish
	wg.Wait()
	close(results)
	close(quit)

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
