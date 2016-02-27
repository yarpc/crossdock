package main

import (
	"fmt"

	"github.com/yarpc/crossdock/execute"
)

// Output results to the console
func Output(results <-chan execute.Result) {
	passed := true
	for result := range results {
		for _, subResult := range result.SubResults {
			if subResult.Status == execute.Failed {
				passed = false
			}
			fmt.Println(subResult)
		}
	}

	fmt.Println(passed)
}
