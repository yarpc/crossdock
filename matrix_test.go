package main

import (
	"testing"

	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
)

var testCase = TestCase{
	Client: "localhost",
	Arguments: Arguments{
		"behavior": "echo",
		"server":   "localhost",
	},
}

func TestExecuteTestCaseWithNon200ShouldFail(t *testing.T) {
	result := ExecuteCase(testCase)
	assert.Equal(t, Failed, result.Status)
	assert.Equal(t, "Internal Server Error\n", result.Response)
}

func TestExecuteTestCaseWithWrongBodyShouldFail(t *testing.T) {
	result := ExecuteCase(testCase)
	assert.Equal(t, Failed, result.Status)
	assert.Equal(t, "ok - not agreed format\n", result.Response)
}
