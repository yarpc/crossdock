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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	plan := New(&Config{
		WaitForHosts: []string{"alpha", "omega"},
		Axes: map[string]Axis{
			"client":    {Name: "client", Values: []string{"alpha", "omega"}},
			"server":    {Name: "server", Values: []string{"alpha", "omega"}},
			"transport": {Name: "transport", Values: []string{"http", "tchannel"}},
		},
		Behaviors: map[string]Behavior{
			"dance": {Name: "dance",
				Clients: "client",
				Params:  []string{"server", "transport"},
			},
			"sing": {Name: "sing",
				Clients: "client",
				Params:  []string{"server", "transport"},
			}}})

	plan.Sort(func(i, j int) bool {
		if plan.TestCases[i].Client != plan.TestCases[j].Client {
			return plan.TestCases[i].Client < plan.TestCases[j].Client
		}
		if plan.TestCases[i].Arguments["server"] != plan.TestCases[j].Arguments["server"] {
			return plan.TestCases[i].Arguments["server"] < plan.TestCases[j].Arguments["server"]
		}
		if plan.TestCases[i].Arguments["transport"] != plan.TestCases[j].Arguments["transport"] {
			return plan.TestCases[i].Arguments["transport"] < plan.TestCases[j].Arguments["transport"]
		}

		return plan.TestCases[i].Arguments["behavior"] < plan.TestCases[j].Arguments["behavior"]
	})

	assert.Equal(t,
		[]TestCase{
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
		},
		plan.TestCases,
	)
}
