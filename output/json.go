package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/crossdock/crossdock/execute"
	"github.com/crossdock/crossdock/plan"
)

type JSONTestReport struct {
	Client    string              `json:"client"`
	Arguments plan.TestClientArgs `json:"arguments"`
	Status    execute.Status      `json:"status"`
	Output    string              `json:"output"`
}

type JSONBehaviorReport struct {
	Params []string         `json:"params"`
	Tests  []JSONTestReport `json:"tests"`
}

type JSONReport struct {
	Behaviors map[string]*JSONBehaviorReport `json:"behaviors"`
}

type JSON struct {
	report JSONReport
	path   string
}

func (j *JSON) Start(plan *plan.Plan) error {
	config := plan.Config
	j.path = config.JSONReportPath
	j.report.Behaviors = make(map[string]*JSONBehaviorReport)

	if config.JSONReportPath == "" {
		return fmt.Errorf("JSON_REPORT_PATH is a required environment variable for REPORT=json")
	}

	for _, behavior := range config.Behaviors {
		behaviorReport := &JSONBehaviorReport{
			Tests:  make([]JSONTestReport, 0, 10),
			Params: behavior.ParamsAxes,
		}
		j.report.Behaviors[behavior.Name] = behaviorReport
	}

	return nil
}

func (j *JSON) Next(test execute.TestResponse) {
	client := test.TestCase.Client
	args := test.TestCase.Arguments
	behavior := test.TestCase.Arguments["behavior"]
	delete(args, "behavior")
	behaviorReport := j.report.Behaviors[behavior]
	if behaviorReport == nil {
		return
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

func (j *JSON) End() error {
	data, err := json.Marshal(j.report)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(j.path, data, 0644); err != nil {
		return err
	}

	return nil
}
