package plan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := Config{
		Clients: []string{"alpha", "omega"},
		Axes: []Axis{{
			Name:   "behavior",
			Values: []string{"dance", "sing"},
		}},
	}
	plan := New(config)

	assert.Equal(t, config, plan.Config)

	assert.Equal(t,
		[]TestCase{
			{
				Client: "alpha",
				Arguments: Arguments{
					"behavior": "dance",
				},
			},
			{
				Client: "alpha",
				Arguments: Arguments{
					"behavior": "sing",
				},
			},
			{
				Client: "omega",
				Arguments: Arguments{
					"behavior": "dance",
				},
			},
			{
				Client: "omega",
				Arguments: Arguments{
					"behavior": "sing",
				},
			},
		},
		plan.TestCases,
	)
}
