package execute

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/yarpc/crossdock/plan"
)

// Run the test program for a given Plan
func Run(plan plan.Plan) <-chan Result {
	results := make(chan Result, 100)
	go func() {
		for _, c := range plan.TestCases {
			results <- executeTestCase(c)
		}
		close(results)
	}()
	return results
}

func executeTestCase(testCase plan.TestCase) Result {
	response, err := makeRequest(testCase)
	if err != nil {
		return Result{
			TestCase: testCase,
			SubResults: []SubResult{{
				Status: Failed,
				Output: fmt.Sprintf("err: %v", err),
			}},
		}
	}

	var subResponses []subResponse
	if err := json.Unmarshal([]byte(response), &subResponses); err != nil {
		return Result{
			TestCase: testCase,
			SubResults: []SubResult{{
				Status: Failed,
				Output: fmt.Sprintf("err: %v", err),
			}},
		}
	}

	return Result{
		TestCase:   testCase,
		SubResults: toSubResults(subResponses),
	}
}

func makeRequest(testCase plan.TestCase) (string, error) {
	callURL, err := url.Parse(fmt.Sprintf("http://%v:8080/", testCase.Client))
	if err != nil {
		return "", err
	}

	args := url.Values{}
	for k, v := range testCase.Arguments {
		args.Add(k, v)
	}
	callURL.RawQuery = args.Encode()

	resp, err := http.Get(callURL.String())
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("wanted status code 200, got %v", resp.StatusCode)
	}

	return string(body), nil
}

type subResponse struct {
	Status string
	Output string
}

func toSubResults(subResponses []subResponse) []SubResult {
	var subResults []SubResult
	for _, subResponse := range subResponses {
		status := Failed
		switch subResponse.Status {
		case "passed":
			status = Success
		case "skipped":
			status = Skipped
		}
		subResult := SubResult{
			Status: status,
			Output: subResponse.Output,
		}
		subResults = append(subResults, subResult)
	}
	return subResults
}
