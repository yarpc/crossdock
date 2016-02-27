package main

import (
	"fmt"
	"os"

	"github.com/yarpc/crossdock/execute"
)

// Output results to the console, if any tests fail exit with non-0
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
	if passed == false {
		fmt.Println("Test suite failed!")
		os.Exit(1)
	}
}
