package plan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromEnviron(t *testing.T) {
	os.Setenv("CROSSDOCK_CLIENTS", "yarpc-go,yarpc-node,yarpc-browser")
	os.Setenv("CROSSDOCK_AXIS_SERVER", "yarpc-go,yarpc-node")
	os.Setenv("CROSSDOCK_AXIS_TRANSPORT", "http,tchannel")

	plan := ReadFromEnviron()

	assert.Equal(t, plan.Clients, []string{"yarpc-go", "yarpc-node", "yarpc-browser"})
	assert.Equal(t, plan.Axes, []Axis{
		{Name: "server", Values: []string{"yarpc-go", "yarpc-node"}},
		{Name: "transport", Values: []string{"http", "tchannel"}},
	})
}
