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
	"sort"
	"strings"

	"github.com/crossdock/crossdock/execute"
	"github.com/crossdock/crossdock/plan"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var grey = color.New(color.FgBlack, color.Bold).SprintFunc()

// List is the output reporter as a list, with compact option.
type List struct {
	Compact bool
	dotted  bool
}

func (list *List) printf(format string, a ...interface{}) (n int, err error) {
	if list.dotted {
		list.dotted = false
		fmt.Println()
	}
	return fmt.Printf(format, a...)
}

func (list *List) printDot(status execute.Status) {
	list.dotted = true
	var dot string
	switch status {
	case execute.Skipped:
		dot = grey("S")
	default:
		dot = "."
	}
	fmt.Printf(dot)
}

func (list *List) Start(plan *plan.Plan) error {
	return nil
}

func statusToColoredSymbol(status execute.Status) string {
	switch status {
	case execute.Success:
		return green("✓")
	case execute.Skipped:
		return grey("S")
	default:
		return red("✗")
	}
}

func fmtTestCase(test plan.TestCase) string {
	var argsList []string
	var behavior string
	for k, v := range test.Arguments {
		if k == "behavior" {
			behavior = v
			continue
		}
		if k == "client" {
			continue
		}
		argsList = append(argsList, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(argsList)
	return fmt.Sprintf("[%v] %v→ (%v)", behavior, test.Client,
		strings.Join(argsList, " "))
}

func (list *List) Next(test execute.TestResponse) {
	if len(test.Results) == 0 {
		list.printf("%v %v ⇒ (no results)\n", blue("∅"), fmtTestCase(test.TestCase))
		return
	}
	if len(test.Results) == 1 {
		result := test.Results[0]
		if list.Compact && result.Status != execute.Failed {
			list.printDot(result.Status)
		} else {
			list.printf("%v %v ⇒ %v\n", statusToColoredSymbol(result.Status),
				fmtTestCase(test.TestCase), result.Output)
		}
		return
	}

	total := len(test.Results)
	passed := 0
	skipped := 0
	for _, result := range test.Results {
		switch result.Status {
		case execute.Success:
			passed++
		case execute.Skipped:
			skipped++
		}
	}
	failed := total - passed - skipped

	behaviorStatus := execute.Failed
	if skipped == total {
		behaviorStatus = execute.Skipped
	} else if failed == 0 {
		behaviorStatus = execute.Success
	}

	if list.Compact && behaviorStatus != execute.Failed {
		list.printDot(behaviorStatus)
		return
	}

	list.printf("%v %v (%v/%v passed, %v/%v skipped)\n",
		statusToColoredSymbol(behaviorStatus), fmtTestCase(test.TestCase),
		passed, total-skipped, skipped, total)

	for idx, result := range test.Results {
		if list.Compact && result.Status != execute.Failed {
			list.printDot(result.Status)
			continue
		}
		prefix := "├"
		if idx == total-1 {
			prefix = "└"
		}
		list.printf("   %v %v ⇒ %v\n", prefix,
			statusToColoredSymbol(result.Status), result.Output)
	}
}

func (list *List) End() error {
	list.printf("")
	return nil
}
