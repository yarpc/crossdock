package output

import (
	"github.com/crossdock/crossdock/execute"
	"github.com/crossdock/crossdock/plan"
)

// Mux multiplexes to multiples reporters.
type Mux []Reporter

func (r Mux) Start(plan *plan.Plan) error {
	for _, reporter := range r {
		if err := reporter.Start(plan); err != nil {
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
