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
