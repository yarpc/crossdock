package execute

import "github.com/yarpc/crossdock/plan"

// FinalResult summarizes the entire test run
type FinalResult struct {
	Passed  bool
	Results []Result
}

// Result contains replies from test clients
type Result struct {
	TestCase plan.TestCase
	Status   Status
	Output   string
}

// Status is an enum that represents test success/failure
type Status int

const (
	// Success indicates a client's TestCase passed
	Success Status = 1 + iota

	// Failed indicates a client's TestCase did not pass
	Failed

	// Skipped indicates a client' TestCase did not run
	Skipped
)
