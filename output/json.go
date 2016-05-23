package output

import (
	"encoding/json"
	"fmt"
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

type JSON struct {
	s Summary
	r JSONReport
}

func (j *JSON) Start(config *plan.Config) error {
	j.r.Behaviors = make(map[string]*JSONBehaviorReport)

	if config.JSONReportPath == "" {
		return fmt.Errorf("JSON_REPORT_PATH is a required environment varialbe for REPORT=json")
	}

	for _, behavior := range config.Behaviors {
		behaviorReport := &JSONBehaviorReport{
			Tests:  make([]JSONTestReport, 0, 10),
			Params: behavior.Params,
		}
		j.r.Behaviors[behavior.Name] = behaviorReport
	}

	return nil
}

func (j *JSON) Next(test execute.TestResponse, config *plan.Config) {
	client := test.TestCase.Client
	args := test.TestCase.Arguments
	behavior := test.TestCase.Arguments["behavior"]
	delete(args, "behavior")
	behaviorReport := j.r.Behaviors[behavior]
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

func (j *JSON) End(config *plan.Config) (Summary, error) {
	data, err := json.Marshal(j.r)
	if err != nil {
		return j.s, err
	}

	err = ioutil.WriteFile(config.JSONReportPath, data, 0644)
	if err != nil {
		return j.s, err
	}

	return j.s, nil
}
