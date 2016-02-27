package output

import (
	"errors"
	"fmt"

	"github.com/yarpc/crossdock/execute"
)

// Stream results to the console, error at end if any fail
func Stream(results <-chan execute.Result) error {
	failed := false
	for result := range results {
		for _, subResult := range result.SubResults {
			var statStr string
			switch subResult.Status {
			case execute.Success:
				statStr = "PASSED"
			case execute.Skipped:
				statStr = "SKIPPED"
			default:
				statStr = "FAILED"
				failed = true
			}
			fmt.Printf("%v - %v - %v\n", statStr, result.TestCase, subResult.Output)
		}
	}
	if failed == true {
		return errors.New("one or more tests failed")
	}
	return nil
}
