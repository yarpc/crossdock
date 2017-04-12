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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/crossdock/crossdock/plan"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// Run the test program for a given Plan
func Run(plan *plan.Plan) <-chan TestResponse {
	tests := make(chan TestResponse, 100)
	go func() {
		for _, c := range plan.TestCases {
			tests <- executeTestCase(c)
		}
		close(tests)
	}()
	return tests
}

func executeTestCase(testCase plan.TestCase) TestResponse {
	if testCase.Skip {
		return TestResponse{
			TestCase: testCase,
			Results: []Result{{
				Status: Skipped,
				Output: fmt.Sprintf("Test filtered out"),
			}},
		}
	}
	response, err := makeRequest(testCase)
	if err != nil {
		return TestResponse{
			TestCase: testCase,
			Results: []Result{{
				Status: Failed,
				Output: fmt.Sprintf("err: %v", err),
			}},
		}
	}

	var subResponses []subResponse
	if err := json.Unmarshal(response, &subResponses); err != nil {
		return TestResponse{
			TestCase: testCase,
			Results: []Result{{
				Status: Failed,
				Output: fmt.Sprintf("err: %v", err),
			}},
		}
	}

	if len(subResponses) == 0 {
		return TestResponse{
			TestCase: testCase,
			Results: []Result{{
				Status: Failed,
				Output: "client returned 0 results",
			}},
		}
	}

	return TestResponse{
		TestCase: testCase,
		Results:  toResults(subResponses),
	}
}

func makeRequest(testCase plan.TestCase) ([]byte, error) {
	callURL, err := url.Parse(fmt.Sprintf("http://%v:8080/", testCase.Client))
	if err != nil {
		return []byte(""), err
	}

	args := url.Values{}
	for k, v := range testCase.Arguments {
		args.Add(k, v)
	}
	callURL.RawQuery = args.Encode()

	ctx, _ := context.WithTimeout(
		context.Background(), testCase.Plan.Config.CallTimeout)
	resp, err := ctxhttp.Get(ctx, nil, callURL.String())
	if err != nil {
		return []byte(""), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	if resp.StatusCode != 200 {
		return []byte(""), fmt.Errorf("wanted status code 200, got %v", resp.StatusCode)
	}

	return body, nil
}

type subResponse struct {
	Status string
	Output string
}

func toResults(subResponses []subResponse) []Result {
	var results []Result
	for _, subResponse := range subResponses {
		status := Failed
		switch subResponse.Status {
		case "passed":
			status = Success
		case "skipped":
			status = Skipped
		}
		result := Result{
			Status: status,
			Output: subResponse.Output,
		}
		results = append(results, result)
	}
	return results
}
