package main

import "github.com/yarpc/crossdock/plan"

// Arguments represents custom args to pass to test client.
type Arguments map[string]string

// TestCase represents the request made to test clients.
type TestCase struct {
	Client    string
	Arguments Arguments
}

// Collect all TestCases for a given Matrix
func Collect(plan plan.Plan) []TestCase {
	var cases []TestCase
	for _, client := range plan.Clients {
		clientTestCases := createCasesForClient(client, plan.Axes, make(map[string]string))
		cases = append(cases, clientTestCases...)
	}

	return cases
}

func createCasesForClient(client string, axes []plan.Axis, args Arguments) []TestCase {
	if len(axes) == 0 {
		return []TestCase{{
			Client:    client,
			Arguments: copyArgs(args),
		}}
	}

	var testCases []TestCase

	axis := axes[0]
	for _, p := range axis.Values {
		args[axis.Name] = p
		testCases = append(testCases, createCasesForClient(client, axes[1:], args)...)
	}
	return testCases
}

func copyArgs(args Arguments) Arguments {
	copy := make(map[string]string)
	for k, v := range args {
		copy[k] = v
	}
	return copy
}
