package plan

import (
	"sort"
	"time"
)

// Config describes the unstructured test plan
type Config struct {
	CallTimeout  time.Duration
	WaitForHosts []string
	Axes         map[string]Axis
	Behaviors    map[string]Behavior
}

// Axis represents combinational args to be passed to the test clients
type Axis struct {
	Name   string
	Values []string
}

// Behavior represents the test behavior will be triggered by crossdock
type Behavior struct {
	Name    string
	Clients string
	Params  []string
}

// Plan describes the entirety of the test program
type Plan struct {
	Config    *Config
	TestCases []TestCase
	less      func(i, j int) bool
}

// Len is part of sort.Interface.
func (p *Plan) Len() int {
	return len(p.TestCases)
}

// Swap is part of sort.Interface.
func (p *Plan) Swap(i, j int) {
	p.TestCases[i], p.TestCases[j] = p.TestCases[j], p.TestCases[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (p *Plan) Less(i, j int) bool {
	return p.less(i, j)
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (p *Plan) Sort(less func(i, j int) bool) {
	p.less = less
	sort.Sort(p)
}

// TestCase represents the request made to test clients.
type TestCase struct {
	Plan      *Plan
	Client    string
	Arguments Arguments
}

// Arguments represents custom args to pass to test client.
type Arguments map[string]string
