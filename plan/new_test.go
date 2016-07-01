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

func TestCombinations(t *testing.T) {
	tests := []struct {
		i [][]string
		o [][]string
	}{
		{
			i: nil,
			o: nil,
		},
		{
			i: [][]string{},
			o: nil,
		},
		{
			i: [][]string{
				{},
			},
			o: nil,
		},
		{
			i: [][]string{
				{"x1"},
				{},
			},
			o: nil,
		},
		{
			i: [][]string{
				{"x1"},
				{"y1"},
			},
			o: [][]string{
				{"x1", "y1"},
			},
		},
		{
			i: [][]string{
				{"x1"},
				{"y1", "y2", "y3"},
			},
			o: [][]string{
				{"x1", "y1"},
				{"x1", "y2"},
				{"x1", "y3"},
			},
		},
		{
			i: [][]string{
				{"x1"},
				{"y1", "y2", "y3"},
				{"z1", "z2"},
			},
			o: [][]string{
				{"x1", "y1", "z1"},
				{"x1", "y2", "z1"},
				{"x1", "y3", "z1"},
				{"x1", "y1", "z2"},
				{"x1", "y2", "z2"},
				{"x1", "y3", "z2"},
			},
		},
		{
			i: [][]string{
				{"x1", "x2"},
				{"y1", "y2"},
				{"z1", "z2"},
			},
			o: [][]string{
				{"x1", "y1", "z1"},
				{"x2", "y1", "z1"},
				{"x1", "y2", "z1"},
				{"x2", "y2", "z1"},
				{"x1", "y1", "z2"},
				{"x2", "y1", "z2"},
				{"x1", "y2", "z2"},
				{"x2", "y2", "z2"},
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.o, combinations(tt.i))
	}
}

func TestNew(t *testing.T) {
	plan := New(&Config{
		WaitForHosts: []string{"alpha", "omega"},
		Axes: Axes{
			{Name: "client", Values: []string{"alpha", "omega"}},
			{Name: "server", Values: []string{"alpha", "omega"}},
			{Name: "transport", Values: []string{"http", "tchannel"}},
		},
		Behaviors: []Behavior{
			{Name: "dance",
				ClientAxis: "client",
				ParamsAxes: []string{"server", "transport"},
			},
			{Name: "sing",
				ClientAxis: "client",
				ParamsAxes: []string{"server", "transport"},
			}}})

	wanted := []TestCase{
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"behavior": "dance",
				"client": "alpha", "server": "alpha", "transport": "http"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"transport": "http", "behavior": "dance",
				"client": "omega", "server": "alpha"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"behavior": "dance", "client": "alpha",
				"server": "omega", "transport": "http"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"behavior": "dance", "client": "omega",
				"server": "omega", "transport": "http"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"behavior": "dance", "client": "alpha",
				"server": "alpha", "transport": "tchannel"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"client": "omega", "server": "alpha",
				"transport": "tchannel", "behavior": "dance"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"server": "omega", "transport": "tchannel",
				"behavior": "dance", "client": "alpha"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"transport": "tchannel", "behavior": "dance",
				"client": "omega", "server": "omega"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"client": "alpha", "server": "alpha",
				"transport": "http", "behavior": "sing"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"behavior": "sing", "client": "omega",
				"server": "alpha", "transport": "http"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"behavior": "sing", "client": "alpha",
				"server": "omega", "transport": "http"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"behavior": "sing", "client": "omega",
				"server": "omega", "transport": "http"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"server": "alpha", "transport": "tchannel",
				"behavior": "sing", "client": "alpha"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"behavior": "sing", "client": "omega",
				"server": "alpha", "transport": "tchannel"}},
		TestCase{Plan: plan, Client: "alpha",
			Arguments: TestClientArgs{"behavior": "sing", "client": "alpha",
				"server": "omega", "transport": "tchannel"}},
		TestCase{Plan: plan, Client: "omega",
			Arguments: TestClientArgs{"behavior": "sing", "client": "omega",
				"server": "omega", "transport": "tchannel"}}}

	assert.Equal(t, plan.TestCases, wanted)
}
