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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFromEnviron(t *testing.T) {
	os.Setenv("REPORT", "list")
	os.Setenv("AXIS_CLIENT", "yarpc-go,yarpc-node,yarpc-browser")
	os.Setenv("AXIS_SERVER", "yarpc-go,yarpc-node")
	os.Setenv("AXIS_TRANSPORT", "http,tchannel")
	os.Setenv("BEHAVIOR_ECHO", "client,server,transport")
	os.Setenv("SKIP_ECHO", "client:yarpc-go+transport:tchannel")
	os.Setenv("CALL_TIMEOUT", "10s")
	os.Setenv("WAIT_FOR_TIMEOUT", "20s")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "cross dock configuration is incorrect.")

	client := Axis{Name: "client", Values: []string{
		"yarpc-browser",
		"yarpc-go",
		"yarpc-node",
	}}
	server := Axis{Name: "server", Values: []string{
		"yarpc-go",
		"yarpc-node",
	}}
	transport := Axis{Name: "transport", Values: []string{
		"http",
		"tchannel",
	}}

	assert.Equal(t, []string{"list"}, config.Reports)

	assert.Equal(t, Axes{client, server, transport}, config.Axes)

	assert.Equal(t, Behaviors{
		{
			Name:       "echo",
			ClientAxis: "client",
			ParamsAxes: []string{"server", "transport"},
			Filters: []Filter{
				Filter{
					AxisMatches: []AxisMatch{
						AxisMatch{
							Name:  "client",
							Value: "yarpc-go",
						},
						AxisMatch{
							Name:  "transport",
							Value: "tchannel",
						},
					},
				},
			},
		},
	}, config.Behaviors)

	assert.Equal(t, 10*time.Second, config.CallTimeout)

	assert.Equal(t, 20*time.Second, config.WaitForTimeout)
}

func TestReadConfigFromEnvironTrimsWhitespace(t *testing.T) {
	os.Setenv("WAIT_FOR", " alpha, omega ")
	defer os.Clearenv()

	config, err := ReadConfigFromEnviron()
	assert.NoError(t, err, "crossdock configuration is incorrect")

	assert.Equal(t, config.WaitForHosts, []string{"alpha", "omega"})
}

func TestParseBehavior(t *testing.T) {
	tests := []struct {
		give string
		want Behavior
	}{
		{
			give: "foo=client,server",
			want: Behavior{
				Name:       "foo",
				ClientAxis: "client",
				ParamsAxes: []string{"server"},
			},
		},
		{
			give: "x=a,b,c,d",
			want: Behavior{
				Name:       "x",
				ClientAxis: "a",
				ParamsAxes: []string{"b", "c", "d"},
			},
		},
		{
			give: "y=c,b,a",
			want: Behavior{
				Name:       "y",
				ClientAxis: "c",
				ParamsAxes: []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, parseBehavior(tt.give))
	}
}

func TestParseSkipBehavior(t *testing.T) {
	tests := []struct {
		give             string
		desc             string
		wantFilters      []Filter
		wantBehaviorName string
		wantError        error
	}{
		{
			give: "foo=client:c+server:b",
			desc: "single filter.",
			wantFilters: []Filter{
				{
					AxisMatches: []AxisMatch{
						AxisMatch{
							Name:  "client",
							Value: "c",
						},
						AxisMatch{
							Name:  "server",
							Value: "b",
						},
					},
				},
			},
			wantBehaviorName: "foo",
			wantError:        nil,
		},
		{
			give: "x=a:b,c:d",
			desc: "multiple filters.",
			wantFilters: []Filter{
				{
					AxisMatches: []AxisMatch{
						AxisMatch{
							Name:  "a",
							Value: "b",
						},
					},
				},
				{
					AxisMatches: []AxisMatch{
						AxisMatch{
							Name:  "c",
							Value: "d",
						},
					},
				},
			},
			wantBehaviorName: "x",
			wantError:        nil,
		},
		{
			give: "foo=client:weird:name+server:b",
			desc: "single filter testing name separation.",
			wantFilters: []Filter{
				{
					AxisMatches: []AxisMatch{
						AxisMatch{
							Name:  "client",
							Value: "weird:name",
						},
						AxisMatch{
							Name:  "server",
							Value: "b",
						},
					},
				},
			},
			wantBehaviorName: "foo",
			wantError:        nil,
		},
		{
			give:             "x=a:b,,c:d",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \"\" in input \"x=a:b,,c:d\" is not of form 'key:value'"),
		},
		{
			give:             "x",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("missing '=' in the input: \"x\""),
		},
		{
			give:             "  x    ",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("missing '=' in the input: \"  x    \""),
		},
		{
			give:             "x=:",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \":\" in input \"x=:\" is empty"),
		},
		{
			give:             "x=",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \"\" in input \"x=\" is not of form 'key:value'"),
		},
		{
			give:             "x  =  ",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \"  \" in input \"x  =  \" is not of form 'key:value'"),
		},
		{
			give:             "x= :",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \" :\" in input \"x= :\" is empty"),
		},
		{
			give:             "x=  :  ",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \"  :  \" in input \"x=  :  \" is empty"),
		},
		{
			give:             "x=     :",
			desc:             "invalid filters",
			wantFilters:      nil,
			wantBehaviorName: "",
			wantError:        fmt.Errorf("match \"     :\" in input \"x=     :\" is empty"),
		},
	}

	for _, tt := range tests {
		behaviorName, filters, err := parseSkipBehavior(tt.give)
		if tt.wantError != nil {
			if assert.Error(t, err) {
				assert.Equal(t, tt.wantError, err)
			}
		} else {
			if assert.NoError(t, err) {
				assert.Equal(t, tt.wantFilters, filters, tt.desc)
				assert.Equal(t, tt.wantBehaviorName, behaviorName, tt.desc)
			}
		}

	}
}
