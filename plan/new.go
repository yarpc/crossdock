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
	for _, behavior := range plan.Config.Behaviors {
		for _, client := range plan.Config.Axes[behavior.Clients].Values {
			combos := recurseCombinations(0, plan, client, behavior, map[string]string{"behavior": behavior.Name})
			testCases = append(testCases, combos...)
		}
	}
	return testCases
}

func recurseCombinations(level int, plan *Plan, client string, behavior Behavior, args Arguments) []TestCase {
	if level == len(behavior.Params) {
		return []TestCase{{
			Plan:      plan,
			Client:    client,
			Arguments: copyArgs(args),
		}}
	}
	var testCases []TestCase
	param := behavior.Params[level]
	for _, axis := range plan.Config.Axes[param].Values {
		args[param] = axis
		testCases = append(testCases, recurseCombinations(level+1, plan, client, behavior, args)...)
	}
	return testCases
}

func copyArgs(args Arguments) Arguments {
	copy := make(map[string]string)
	for k, v := range args {
		copy[k] = v
	}
	return copy
}
