package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Wait for hosts to become ready
func Wait(hosts []string, timeout time.Duration) {
	begin := time.Now()
	var wg sync.WaitGroup
	cancel := make(chan struct{})

	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			WaitForHTTPRequest(host, cancel)
			wg.Done()
		}(fmt.Sprintf("%s:8080", host))
	}

	timer := time.AfterFunc(timeout, func() {
		close(cancel)
	})

	wg.Wait()

	if !timer.Stop() {
		log.Fatalf("Error: One or more services timed out after %d second(s)", timeout)
	}
	fmt.Printf("\nAll services are up after %v!\n", time.Since(begin))
}

// WaitForHTTPRequest polls host until it can make a request
func WaitForHTTPRequest(host string, cancel <-chan struct{}) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/",
	}

	err := errors.New("init")
	for err != nil {
		req, reqErr := http.NewRequest("HEAD", url.String(), nil)
		if reqErr != nil {
			log.Printf("Warning: Failed to create request for URL '%s' -  skipping service '%s'",
				url.String(), host)
			return
		}

		tr := &http.Transport{}
		client := &http.Client{Transport: tr}
		c := make(chan error, 1)
		go func() {
			_, err := client.Do(req)
			c <- err
		}()

		select {
		case <-cancel:
			tr.CancelRequest(req)
			log.Printf("HTTP: Service %v timed out. Last error: %v", host, err)
			<-c
			return
		case err = <-c:
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("HTTP: Service %v is up\n", host)
}
