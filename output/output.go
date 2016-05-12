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

package output

import (
	"fmt"

	"github.com/yarpc/crossdock/execute"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

// Summary contains an account of the test run
type Summary struct {
	Failed bool

	SuccessAmount int
	FailAmount    int
	SkippedAmount int
}

// Stream results to the console, error at end if any fail
func Stream(tests <-chan execute.TestResponse) Summary {
	var summary Summary
	for test := range tests {
		for _, result := range test.Results {
			var statStr string
			switch result.Status {
			case execute.Success:
				statStr = green("✓")
				summary.SuccessAmount++
			case execute.Skipped:
				statStr = yellow("S")
				summary.SkippedAmount++
			default:
				statStr = red("F")
				summary.Failed = true
				summary.FailAmount++
			}
			fmt.Printf("%v - %v - %v\n", statStr, test.TestCase, result.Output)
		}
	}
	return summary
}

// Summarize outputs the summary to the console
func Summarize(summary Summary) {
	fmt.Println("")
	if summary.SuccessAmount > 0 {
		fmt.Printf("%v passed\n", summary.SuccessAmount)
	}
	if summary.FailAmount > 0 {
		fmt.Printf("%v failed\n", summary.FailAmount)
	}
	if summary.SkippedAmount > 0 {
		fmt.Printf("%v skipped\n", summary.SkippedAmount)
	}

	if summary.Failed == true {
		fmt.Printf("\nTests did not pass!\n\n")
		return
	}
	fmt.Printf("\nTests passed!\n\n")
}
