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

	"github.com/crossdock/crossdock/execute"
	"github.com/crossdock/crossdock/plan"
)

// Compact is a simple Reporter that prints the test output to stdout with lots of dots (".').
type Compact struct {
	dotted bool
}

func (compact *Compact) printf(format string, a ...interface{}) (n int, err error) {
	if compact.dotted {
		compact.dotted = false
		fmt.Println()
	}
	return fmt.Printf(format, a...)
}

func (compact *Compact) printDot(status execute.Status) {
	compact.dotted = true
	var dot string
	switch status {
	case execute.Skipped:
		dot = grey("S")
	default:
		dot = "."
	}
	fmt.Printf(dot)
}

func (compact *Compact) Start(plan *plan.Plan) error {
	return nil
}

func (compact *Compact) Next(test execute.TestResponse) {
	if len(test.Results) == 0 {
		compact.printf("%v %v ⇒ (no results)\n", blue("∅"), fmtTestCase(test.TestCase))
		return
	}

	if len(test.Results) == 1 {
		result := test.Results[0]
		if result.Status != execute.Failed {
			compact.printDot(result.Status)
		} else {
			compact.printf("%v %v ⇒ %v\n", statusToColoredSymbol(result.Status),
				fmtTestCase(test.TestCase), result.Output)
		}
		return
	}

	bs := test.SummarizeStatus()

	if bs.Status != execute.Failed {
		compact.printDot(bs.Status)
		return
	}

	fmt.Printf("%v %v (%v/%v passed, %v/%v skipped)\n",
		statusToColoredSymbol(bs.Status), fmtTestCase(test.TestCase),
		bs.Passed, bs.Total-bs.Skipped, bs.Skipped, bs.Total)

	for idx, result := range test.Results {
		if result.Status != execute.Failed {
			compact.printDot(result.Status)
			continue
		}
		prefix := "├"
		if idx == bs.Total-1 {
			prefix = "└"
		}
		compact.printf("   %v %v ⇒ %v\n", prefix,
			statusToColoredSymbol(result.Status), result.Output)
	}
}

func (compact *Compact) End() error {
	compact.printf("")
	return nil
}
