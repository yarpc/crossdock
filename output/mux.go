package output

import (
	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
)

type Mux struct {
	s         Summary
	Reporters []Reporter
}

func (r Mux) Start(config *plan.Config) error {
	for _, reporter := range r.Reporters {
		err := reporter.Start(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Mux) Next(response execute.TestResponse, config *plan.Config) {
	for _, reporter := range r.Reporters {
		reporter.Next(response, config)
	}
}

func (r Mux) End(config *plan.Config) (Summary, error) {
	var summary Summary
	for _, reporter := range r.Reporters {
		summary, err := reporter.End(config)
		if err != nil {
			return summary, err
		}
	}
	return summary, nil
}
