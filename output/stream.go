package output

import (
	"errors"
	"fmt"

	"github.com/yarpc/crossdock/execute"
)

// Stream results to the console, error at end if any fail
func Stream(results <-chan execute.Result) error {
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
		return errors.New("one or more tests failed")
	}
	return nil
}
