package output

import (
	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
)

type Mux []Reporter

func (r Mux) Start(config *plan.Config) error {
	for _, reporter := range r {
		if err := reporter.Start(config); err != nil {
			return err
		}
	}
	return nil
}

func (r Mux) Next(response execute.TestResponse) {
	for _, reporter := range r {
		reporter.Next(response)
	}
}

func (r Mux) End() error {
	for _, reporter := range r {
		if err := reporter.End(); err != nil {
			return err
		}
	}
	return nil
}
