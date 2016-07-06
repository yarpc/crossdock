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

package execute

import (
	"fmt"

	"github.com/crossdock/crossdock/plan"
)

// TestResponse contains the reply from a test client, most importantly,
// contains a list of Results for the test cases ran by the test client
type TestResponse struct {
	TestCase plan.TestCase
	Results  []Result
}

// StatusSummary contains the status report of a TestReponse, see
// SummarizeStatus.
type StatusSummary struct {
	Status  Status
	Total   int
	Skipped int
	Passed  int
	Failed  int
}

// SummarizeStatus returns a summary of the TestResponse, compiling it from all
// the Results pertaining to it.
func (test *TestResponse) SummarizeStatus() StatusSummary {
	total := len(test.Results)
	passed := 0
	skipped := 0
	for _, result := range test.Results {
		switch result.Status {
		case Success:
			passed++
		case Skipped:
			skipped++
		}
	}
	failed := total - passed - skipped

	behaviorStatus := Failed
	if skipped == total {
		behaviorStatus = Skipped
	} else if failed == 0 {
		behaviorStatus = Success
	}

	return StatusSummary{
		Status:  behaviorStatus,
		Total:   total,
		Skipped: skipped,
		Passed:  passed,
		Failed:  failed,
	}
}

// Result is the outcome of an individual test case ran by the test client
type Result struct {
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

func (s Status) String() string {
	switch s {
	case Success:
		return "success"
	case Failed:
		return "failed"
	case Skipped:
		return "skipped"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}
