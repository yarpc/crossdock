# Write Test Client

A test client is just a HTTP server written in any language, that listens on port `8080`.
Crossdock executes the test matrix against each client.

The following illustrates an over-simplified test client written in Go (available in the `example/client.go`):

```go
package main

import "net/http"

func main() {
	// Crossdock makes all calls to http://<test-client>:8080/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// custom arguments, called axis, are configured in
		// docker-compose.yml and then passed as query params like so:
		// http://<test-client>:8080/?behavior=dance
		behavior := r.FormValue("behavior")

		// when client is called with no arguments,
		// report back with a 200 when ready to run tests
		if behavior == "" {
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
```

[Run Crossdock →](run-crossdock.md)
