package plan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFromEnviron(t *testing.T) {
	os.Setenv("CROSSDOCK_AXIS_CLIENT", "yarpc-go,yarpc-node,yarpc-browser")
	os.Setenv("CROSSDOCK_AXIS_SERVER", "yarpc-go,yarpc-node")
	os.Setenv("CROSSDOCK_AXIS_TRANSPORT", "http,tchannel")
	os.Setenv("CROSSDOCK_BEHAVIOR_ECHO", "client,server,transport")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "cross dock configuration is incorrect.")

	client := Axis{Name: "client", Values: []string{"yarpc-go", "yarpc-node", "yarpc-browser"}}
	server := Axis{Name: "server", Values: []string{"yarpc-go", "yarpc-node"}}
	transport := Axis{Name: "transport", Values: []string{"http", "tchannel"}}

	assert.Equal(t, config.Axes, map[string]Axis{
		"client":    client,
		"server":    server,
		"transport": transport,
	})

	assert.Equal(t, config.Behaviors, map[string]Behavior{
		"echo": {
			Name:    "echo",
			Clients: "client",
			Params:  []string{"server", "transport"},
		}})
}

func TestReadConfigFromEnvironTrimsWhitespace(t *testing.T) {
	os.Setenv("CROSSDOCK_WAIT_FOR", " alpha, omega ")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "crossdock configuration is incorrect")

	assert.Equal(t, config.WaitForHosts, []string{"alpha", "omega"})
}
