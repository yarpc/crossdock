package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Matrix struct {
	Clients   []string
	Servers   []string
	Behaviors []string
}

type Results []Result

func BeginMatrixTest(matrix Matrix) Results {

	len := len(matrix.Clients) * len(matrix.Servers) * len(matrix.Behaviors)
	results := make(Results, 0, len)

	for _, client := range matrix.Clients {
		for _, server := range matrix.Servers {
			for _, behavior := range matrix.Behaviors {

				testCase := TestCase{
					Client:   client,
					Server:   server,
					Behavior: behavior,
				}

				result := ExecuteTestCase(testCase)
				results = append(results, result)
			}
		}
	}

	return results
}

type TestCase struct {
	Client   string
	Server   string
	Behavior string
}

type Result struct {
	TestCase TestCase
	Status   Status
	Response string
}

type Status int

const (
	Success Status = 1 + iota
	Failed
	Skipped
)

func ExecuteTestCase(testCase TestCase) Result {

	callUrl, err := url.Parse(fmt.Sprintf("http://%v:8080/", testCase.Client))
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	params.Add("server", testCase.Server)
	params.Add("behavior", testCase.Behavior)
	callUrl.RawQuery = params.Encode()

	// Produce TAP header
	fmt.Printf("# %s %s %s\n", testCase.Client, testCase.Server, testCase.Behavior)

	resp, err := http.Get(callUrl.String())
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	sbody := string(body)
	status := Skipped
	if resp.StatusCode != 200 {
		status = Failed
		fmt.Printf("not ok - harness responded with HTTP error %s\n", resp.StatusCode)
	} else if strings.HasPrefix(sbody, "not ok ") {
		status = Failed
		fmt.Printf("%s\n", sbody)
	} else if strings.HasPrefix(sbody, "ok ") {
		status = Success
		fmt.Printf("%s\n", sbody)
	} else {
		fmt.Printf("%s\n", sbody)
	}

	return Result{
		TestCase: testCase,
		Status:   status,
		Response: sbody,
	}
}
