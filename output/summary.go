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

package output

import (
	"errors"
	"fmt"

	"github.com/crossdock/crossdock/execute"
	"github.com/crossdock/crossdock/plan"
)

// Summary contains an account of the test run
type Summary struct {
	Failed bool

	NumSuccess int
	NumFail    int
	NumSkipped int
}

func (s *Summary) Start(config *plan.Config) error {
	return nil
}

func (s *Summary) Next(test execute.TestResponse) {
	for _, result := range test.Results {
		switch result.Status {
		case execute.Success:
			s.NumSuccess++
		case execute.Skipped:
			s.NumSkipped++
		default:
			s.Failed = true
			s.NumFail++
		}
	}
}

func (s *Summary) End() error {
	fmt.Println("")
	if s.NumSuccess > 0 {
		fmt.Printf("%v passed\n", s.NumSuccess)
	}
	if s.NumFail > 0 {
		fmt.Printf("%v failed\n", s.NumFail)
	}
	if s.NumSkipped > 0 {
		fmt.Printf("%v skipped\n", s.NumSkipped)
	}

	if s.Failed {
		fmt.Printf("\nTests did not pass!\n\n")
		return errors.New("At least one test failed")
	}
	fmt.Printf("\nTests passed!\n\n")
	return nil
}
