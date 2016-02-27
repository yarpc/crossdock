package execute

import "github.com/yarpc/crossdock/plan"

// Result contains total reply from test client
type Result struct {
	TestCase   plan.TestCase
	SubResults []SubResult
}

// SubResult contains result from test client
type SubResult struct {
	Status Status
	Output string
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
