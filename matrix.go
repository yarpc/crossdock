package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Axis represents combinational args to be passed to the test clients
type Axis struct {
	Name   string
	Values []string
}

// Matrix describes the entirety of the test program
type Matrix struct {
	Clients []string
	Axes    []Axis
}

// Result contains replies from test clients
type Result struct {
	TestCase TestCase
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

// Execute the test program for a given Matrix
func Execute(matrix Matrix) []Result {
	cases := Collect(matrix)
	results := ExecuteCases(cases)
	return results
}

// ExecuteCases runs a list of TestCase's
func ExecuteCases(cases []TestCase) []Result {
	var results []Result
	for _, c := range cases {
		results = append(results, ExecuteCase(c)...)
	}
	return results
}

// ExecuteCase fires an HTTP request to the test client
func ExecuteCase(c TestCase) []Result {
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
