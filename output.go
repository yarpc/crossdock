package main

import (
	"fmt"

	"github.com/yarpc/crossdock/execute"
)

// Output results to the console
func Output(results <-chan execute.Result) {
	for result := range results {
		for _, subResult := range result.SubResults {
			fmt.Println(subResult)
		}
	}
}
