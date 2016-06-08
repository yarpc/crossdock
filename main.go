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

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/output"
	"github.com/yarpc/crossdock/plan"
)

func main() {
	fmt.Printf("\nCrossdock starting...\n\n")

	config, err := plan.ReadConfigFromEnviron()
	if err != nil {
		log.Fatal(err)
	}
	plan := plan.New(config)

	fmt.Printf("Waiting on WAIT_FOR=%v\n\n", plan.Config.WaitForHosts)
	execute.Wait(plan.Config.WaitForHosts, time.Duration(30)*time.Second)

	fmt.Printf("\nExecuting Matrix...\n\n")
	results := execute.Run(plan)

	fail := func(err error) {
		os.Exit(1)
	}

	reporter, err := output.GetReporter(config.Reports)
	if err != nil {
		fail(err)
	}

	if err = reporter.Start(config); err != nil {
		fail(err)
	}
	for test := range results {
		reporter.Next(test)
	}
	if err := reporter.End(); err != nil {
		fail(err)
	}
}
