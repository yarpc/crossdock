package plan

import "time"

// Config describes the unstructured test plan
type Config struct {
	CallTimeout time.Duration
	Waiters     []string
	Axes        map[string]Axis
	Behaviors   map[string]Behavior
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
}

// TestCase represents the request made to test clients.
type TestCase struct {
	Plan      *Plan
	Client    string
	Arguments Arguments
}

// Arguments represents custom args to pass to test client.
type Arguments map[string]string
