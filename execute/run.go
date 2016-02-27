package execute

import (
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
	result := Result{
		TestCase: testCase,
	}

	response, err := makeRequest(testCase)
	if err != nil {
		result.SubResults = []SubResult{{
			Status: Failed,
			Output: fmt.Sprintf("err: %v", err),
		}}
		return result
	}

	result.SubResults = []SubResult{{
		Status: Success,
		Output: response,
	}}

	return result
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

type jsonResponse struct {
	SubResults []jsonSubResponse
}

type jsonSubResponse struct {
	Status string `json:"status"`
	Output string `json:"output"`
}
