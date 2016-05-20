package output

import (
	"encoding/json"
	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
	"io/ioutil"
)

var StatusNames = map[execute.Status]string{
	execute.Success: "success",
	execute.Failed:  "failed",
	execute.Skipped: "skipped",
}

type JsonTestReport struct {
	Client    string         `json:"client"`
	Arguments plan.Arguments `json:"arguments"`
	Status    string         `json:"status"`
	Output    string         `json:"output"`
}

type JsonBehaviorReport struct {
	Params []string         `json:"params"`
	Tests  []JsonTestReport `json:"tests"`
}

type JsonReport struct {
	Behaviors map[string]*JsonBehaviorReport `json:"behaviors"`
}

var Json ReporterFunc = func(config *plan.Config, tests <-chan execute.TestResponse) (summary Summary) {

	report := JsonReport{}
	report.Behaviors = make(map[string]*JsonBehaviorReport)

	for _, behavior := range config.Behaviors {
		behaviorReport := new(JsonBehaviorReport)
		behaviorReport.Tests = make([]JsonTestReport, 0, 10)
		behaviorReport.Params = behavior.Params
		report.Behaviors[behavior.Name] = behaviorReport
	}

	for test := range tests {
		client := test.TestCase.Client
		args := test.TestCase.Arguments
		behavior := test.TestCase.Arguments["behavior"]
		delete(args, "behavior")
		behaviorReport := report.Behaviors[behavior]
		if behaviorReport == nil {
			continue
		}
		for _, result := range test.Results {
			behaviorReport.Tests = append(behaviorReport.Tests, JsonTestReport{
				Client:    client,
				Arguments: args,
				Status:    StatusNames[result.Status],
				Output:    result.Output,
			})
		}
	}

	data, err := json.Marshal(report)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/crossdock/crossdock.json", data, 0644)
	if err != nil {
		panic(err)
	}

	return
}
