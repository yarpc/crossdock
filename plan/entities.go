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

import "time"

// Config describes the unstructured test plan
type Config struct {
	Reports              []string
	CallTimeout          time.Duration
	WaitForTimeout       time.Duration
	WaitForHosts         []string
	Axes                 Axes
	Behaviors            Behaviors
	JSONReportPath       string
	HorizontalScaleShard  int
	HorizontalScaleCount int
}

// Axes is a collection of Axis objects sortable by axis name.
type Axis struct {
	Name   string
	Values []string
}

// Axes is a slice of "Axis"
type Axes []Axis

func (a Axes) Len() int           { return len(a) }
func (a Axes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Axes) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Index returns the Axes indexed by name of Axis.
func (a Axes) Index() map[string]Axis {
	axes := make(map[string]Axis, len(a))
	for _, axis := range a {
		axes[axis.Name] = axis
	}
	return axes
}

// Behavior represents the test behavior that will be triggered by crossdock
type Behavior struct {
	Name       string
	ClientAxis string
	ParamsAxes []string
}

// Behaviors is a collection of Behavior objects sortable by behavior name.
type Behaviors []Behavior

func (b Behaviors) Len() int           { return len(b) }
func (b Behaviors) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b Behaviors) Less(i, j int) bool { return b[i].Name < b[j].Name }

// Plan describes the entirety of the test program
type Plan struct {
	Config    *Config
	TestCases []TestCase
	less      func(i, j int) bool
}

// TestCase represents the request made to test clients.
type TestCase struct {
	Plan      *Plan
	Client    string
	Arguments TestClientArgs
}

// TestClientArgs represents custom args to pass to test client.
type TestClientArgs map[string]string
