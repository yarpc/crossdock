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

package plan

import (
	"fmt"
	"strings"
)

// New creates a Plan given a Config
func New(config *Config) *Plan {
	plan := &Plan{
		Config: config,
	}
	plan.TestCases = buildTestCases(plan)
	return plan
}

func buildTestCases(plan *Plan) []TestCase {
	var testCases []TestCase
	axesIndex := plan.Config.Axes.Index()
	for _, behavior := range plan.Config.Behaviors {
		var selectedAxes [][]string
		selectedAxes = append(selectedAxes, axesIndex[behavior.ClientAxis].Values)
		for _, paramAxis := range behavior.ParamsAxes {
			selectedAxes = append(selectedAxes, axesIndex[paramAxis].Values)
		}
		for _, combination := range combinations(selectedAxes) {
			client := combination[0]     // first element is the Client.
			arguments := combination[1:] // remaining is test arguments in order.
			testArgs := TestClientArgs{
				"behavior":          behavior.Name,
				behavior.ClientAxis: client,
			}
			for idx, arg := range arguments {
				testArgs[behavior.ParamsAxes[idx]] = arg
			}
			t := TestCase{
				Plan:      plan,
				Client:    client,
				Arguments: testArgs,
				Skip:      false,
			}

			for _, filter := range behavior.Filters {
				if filter.Matches(testArgs) {
					t.Skip = true
					t.SkipReason = fmt.Sprintf("SKIP_%s=%s", strings.ToUpper(behavior.Name), strings.Join(filter.String(), "+"))
					break
				}
			}
			testCases = append(testCases, t)
		}
	}
	return testCases
}

// combinations takes multiple lists of strings and return multiple lists of
// all the possible combinations. Sublists don't need to be of the same size,
// the output order is dependent of the input order.
// eg: [[x1, x2], [y1, y2]] -> [[x1, y1], [x1, y2], [x2, y1], [x2, y2]]
func combinations(pool [][]string) [][]string {
	if len(pool) == 0 {
		return nil
	}
	if len(pool) == 1 {
		var r [][]string
		for _, v := range pool[0] {
			r = append(r, []string{v})
		}
		return r
	}
	var r [][]string
	for _, remaining := range combinations(pool[1:]) {
		for _, head := range pool[0] {
			r = append(r, append([]string{head}, remaining...))
		}
	}
	return r
}
