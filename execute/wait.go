// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package execute

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
			waitForHTTPRequest(host, cancel)
			wg.Done()
		}(fmt.Sprintf("%s:8080", host))
	}

	timer := time.AfterFunc(timeout, func() {
		close(cancel)
	})

	wg.Wait()

	if !timer.Stop() {
		log.Fatalf("Error: One or more services timed out after %v", timeout)
	}
	fmt.Printf("\nAll services are up after %v!\n", time.Since(begin))
}

// WaitForHTTPRequest polls host until it can make a request
func waitForHTTPRequest(host string, cancel <-chan struct{}) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/",
	}

	err := errors.New("init")
	for err != nil {
		log.Printf("HTTP: HEAD %s\n", url.String())
		req, reqErr := http.NewRequest("HEAD", url.String(), nil)
		if reqErr != nil {
			log.Printf("Warning: Failed to create request for URL '%s' -  skipping service '%s'",
				url.String(), host)
			return
		}
		req.Cancel = cancel

		client := &http.Client{Transport: &http.Transport{}, Timeout: 500 * time.Millisecond}
		c := make(chan error, 1)
		go func() {
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode != http.StatusOK {
				err = errors.New(fmt.Sprintf("Expecting %v: got %v",
					http.StatusOK, resp.Status))
			}
			c <- err
		}()

		select {
		case <-cancel:
			log.Printf("HTTP: Service %v timed out. Last error: %v", host, err)
			<-c
			return
		case err = <-c:
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("HTTP: Service %v is up\n", host)
}
