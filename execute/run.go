package execute

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/yarpc/crossdock/plan"
)

// Run the test program for a given Plan
func Run(plan plan.Plan) FinalResult {
	var results []Result
	for _, c := range plan.TestCases {
		results = append(results, executeCase(c)...)
	}
	final := FinalResult{
		Passed:  true,
		Results: results,
	}
	return final
}

// ExecuteCase fires an HTTP request to the test client
func executeCase(c plan.TestCase) []Result {
	callURL, err := url.Parse(fmt.Sprintf("http://%v:8080/", c.Client))
	if err != nil {
		log.Fatal(err)
	}

	args := url.Values{}
	for k, v := range c.Arguments {
		args.Add(k, v)
	}
	callURL.RawQuery = args.Encode()

	resp, err := http.Get(callURL.String())
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	status := Success
	if resp.StatusCode != 200 {
		status = Failed
	}

	return []Result{{
		TestCase: c,
		Status:   status,
		Output:   string(body),
	}}
}
