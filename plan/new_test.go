package plan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	plan := New(&Config{
		WaitForHosts: []string{"alpha", "omega"},
		Axes: map[string]Axis{
			"client":    {Name: "client", Values: []string{"alpha", "omega"}},
			"server":    {Name: "server", Values: []string{"alpha", "omega"}},
			"transport": {Name: "transport", Values: []string{"http", "tchannel"}},
		},
		Behaviors: map[string]Behavior{
			"dance": {Name: "dance",
				Clients: "client",
				Params:  []string{"server", "transport"},
			},
			"sing": {Name: "sing",
				Clients: "client",
				Params:  []string{"server", "transport"},
			}}})

	plan.Sort(func(i, j int) bool {
		if plan.TestCases[i].Client != plan.TestCases[j].Client {
			return plan.TestCases[i].Client < plan.TestCases[j].Client
		}
		if plan.TestCases[i].Arguments["server"] != plan.TestCases[j].Arguments["server"] {
			return plan.TestCases[i].Arguments["server"] < plan.TestCases[j].Arguments["server"]
		}
		if plan.TestCases[i].Arguments["transport"] != plan.TestCases[j].Arguments["transport"] {
			return plan.TestCases[i].Arguments["transport"] < plan.TestCases[j].Arguments["transport"]
		}

		return plan.TestCases[i].Arguments["behavior"] < plan.TestCases[j].Arguments["behavior"]
	})

	assert.Equal(t,
		[]TestCase{
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "alpha",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "alpha",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "http",
					"behavior":  "sing",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "dance",
				},
			},
			{
				Plan:   plan,
				Client: "omega",
				Arguments: Arguments{
					"server":    "omega",
					"transport": "tchannel",
					"behavior":  "sing",
				},
			},
		},
		plan.TestCases,
	)
}
