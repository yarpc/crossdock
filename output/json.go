package output

import (
	"encoding/json"
	"io/ioutil"

	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
)

type JSONTestReport struct {
	Client    string         `json:"client"`
	Arguments plan.Arguments `json:"arguments"`
	Status    execute.Status `json:"status"`
	Output    string         `json:"output"`
}

type JSONBehaviorReport struct {
	Params []string         `json:"params"`
	Tests  []JSONTestReport `json:"tests"`
}

type JSONReport struct {
	Behaviors map[string]*JSONBehaviorReport `json:"behaviors"`
}

var JSON ReporterFunc = func(config *plan.Config, tests <-chan execute.TestResponse) Summary {
	summary := Summary{}
	report := JSONReport{}
	report.Behaviors = make(map[string]*JSONBehaviorReport)

	for _, behavior := range config.Behaviors {
		behaviorReport := &JSONBehaviorReport{
			Tests:  make([]JSONTestReport, 0, 10),
			Params: behavior.Params,
		}
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
			behaviorReport.Tests = append(behaviorReport.Tests, JSONTestReport{
				Client:    client,
				Arguments: args,
				Status:    result.Status,
				Output:    result.Output,
			})
		}
	}

	data, err := json.Marshal(report)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(config.JSONReportPath, data, 0644)
	if err != nil {
		panic(err)
	}

	return summary
}
