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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testHandler struct {
	t             *testing.T
	delay         time.Duration // delay before replying
	notReadyCount int32         // 500 that many time before returning success.
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := h.t
	if h.delay != 0 {
		t.Logf("testHandler: Delaying for %v", h.delay)
		time.Sleep(h.delay)
	}
	if h.notReadyCount > 0 {
		h.notReadyCount--
		t.Logf("testHandler: Returning %v (left: %v)",
			http.StatusServiceUnavailable, h.notReadyCount)
		http.Error(w, "Server not ready", http.StatusServiceUnavailable)
		return
	}
	if r.Method == "HEAD" {
		w.Header().Add("Content-Length", "0")
		t.Logf("testHandler: Returning %v", http.StatusOK)
		return
	}
	http.NotFound(w, r)
}

func TestObvious(t *testing.T) {
	ts := httptest.NewServer(&testHandler{t: t})
	defer ts.Close()
	host := ts.Listener.Addr().String()
	cancel := make(chan struct{})
	timer := time.AfterFunc(2*time.Second, func() {
		close(cancel)
	})
	waitForHTTPRequest(host, cancel)
	require.True(t, timer.Stop(), "test shouldn't take more than 2 seconds")
}

func TestTooSlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ts := httptest.NewServer(&testHandler{t: t, delay: time.Second})
	defer ts.Close()
	host := ts.Listener.Addr().String()
	cancel := make(chan struct{})
	timer := time.AfterFunc(2*time.Second, func() {
		close(cancel)
	})
	waitForHTTPRequest(host, cancel)
	require.False(t, timer.Stop(), "test should timeout")
}

func TestNotReadyThenOk(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ts := httptest.NewServer(&testHandler{t: t, delay: 100 * time.Millisecond, notReadyCount: 3})
	defer ts.Close()
	host := ts.Listener.Addr().String()
	cancel := make(chan struct{})
	timer := time.AfterFunc(2*time.Second, func() {
		close(cancel)
	})
	waitForHTTPRequest(host, cancel)
	require.True(t, timer.Stop(), "test shouldn't take more than 2 seconds")
}
