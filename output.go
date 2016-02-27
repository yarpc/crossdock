package main

import (
	"fmt"

	"github.com/yarpc/crossdock/execute"
)

// Output results to the console
func Output(results <-chan execute.Result) {
	for result := range results {
		fmt.Println(result)
	}
}
