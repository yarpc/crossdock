// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package plan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFromEnviron(t *testing.T) {
	os.Setenv("REPORT", "list")
	os.Setenv("AXIS_CLIENT", "yarpc-go,yarpc-node,yarpc-browser")
	os.Setenv("AXIS_SERVER", "yarpc-go,yarpc-node")
	os.Setenv("AXIS_TRANSPORT", "http,tchannel")
	os.Setenv("BEHAVIOR_ECHO", "client,server,transport")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "cross dock configuration is incorrect.")

	client := Axis{Name: "client", Values: []string{"yarpc-go", "yarpc-node", "yarpc-browser"}}
	server := Axis{Name: "server", Values: []string{"yarpc-go", "yarpc-node"}}
	transport := Axis{Name: "transport", Values: []string{"http", "tchannel"}}

	assert.Equal(t, config.Reports, []string{"list"})

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
	os.Setenv("WAIT_FOR", " alpha, omega ")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "crossdock configuration is incorrect")

	assert.Equal(t, config.WaitForHosts, []string{"alpha", "omega"})
}
