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

package main

import "net/http"

func main() {
	// Crossdock makes all calls to http://<test-client>:8080/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// custom arguments, called axis, are configured in
		// docker-compose.yml and then passed as query params like so:
		// http://<test-client>:8080/?behavior=dance
		behavior := r.FormValue("behavior")

		// when method is HEAD, report back with a 200 when ready to run tests
		if r.Method == "HEAD" {
			w.Header().Add("Content-Length", "0")
			return
		}

		// once the client is ready, Crossdock will make an HTTP request
		// to / with axis you defined in docker-compose.yml,
		// in this case, we've defined AXIS_BEHAVIOR
		switch behavior {

		// we've recieved a request asking us to test the "dance" behavior,
		// write whatever code we need to verify that behavior, then
		// respond in TAP (testanything.org) format: simply "ok" or "not ok"
		case "dance":
			w.Write([]byte("ok\n"))
			return

		case "run":
			// do something to test the "run" behavior...
			w.Write([]byte("ok\n"))
			return

		default:
			// give a 404 when test is not implemented,
			// Crossdock will mark every 404 test case as "skipped"
			http.NotFound(w, r)
		}
	})
	http.ListenAndServe(":8080", nil)
}
