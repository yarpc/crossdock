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
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	defaultCallTimeout    = 5 * time.Second
	defaultWaitForTimeout = 30 * time.Second
)

// ReadConfigFromEnviron creates a Config by looking for environment variables
func ReadConfigFromEnviron() (*Config, error) {
	const (
		reportKey         = "REPORT"
		callTimeoutKey    = "CALL_TIMEOUT"
		waitForTimeoutKey = "WAIT_FOR_TIMEOUT"
		waitKey           = "WAIT_FOR"
		axisKeyPrefix     = "AXIS_"
		behaviorKeyPrefix = "BEHAVIOR_"
		skipKeyPrefix     = "SKIP_"
		jsonReportPathKey = "JSON_REPORT_PATH"
	)

	callTimeout, _ := time.ParseDuration(os.Getenv(callTimeoutKey))
	if callTimeout == 0 {
		callTimeout = defaultCallTimeout
	}
	waitForTimeout, _ := time.ParseDuration(os.Getenv(waitForTimeoutKey))
	if waitForTimeout == 0 {
		waitForTimeout = defaultWaitForTimeout
	}

	waitForHosts := trimCollection(strings.Split(os.Getenv(waitKey), ","))
	reports := trimCollection(strings.Split(strings.ToLower(os.Getenv(reportKey)), ","))

	var axes Axes
	var behaviors Behaviors
	filterMap := make(map[string][]Filter)
	for _, e := range os.Environ() {
		switch {
		case strings.HasPrefix(e, axisKeyPrefix):
			axis := parseAxis(strings.TrimPrefix(e, axisKeyPrefix))
			axes = append(axes, axis)
		case strings.HasPrefix(e, skipKeyPrefix):
			key, filters, err := parseSkipBehavior(strings.TrimPrefix(e, skipKeyPrefix))
			if err != nil {
				return nil, fmt.Errorf("failed to parse filters from %q: %v", e, err)
			}
			filterMap[key] = filters
		case strings.HasPrefix(e, behaviorKeyPrefix):
			behavior := parseBehavior(strings.TrimPrefix(e, behaviorKeyPrefix))
			behaviors = append(behaviors, behavior)
		}
	}
	sort.Sort(axes)
	sort.Sort(behaviors)
	if err := behaviors.validateAndApplyFilters(filterMap); err != nil {
		return nil, fmt.Errorf("failed to validate filters: %v", err)
	}

	jsonReportPath := os.Getenv(jsonReportPathKey)

	config := &Config{
		Reports:        reports,
		CallTimeout:    callTimeout,
		WaitForTimeout: waitForTimeout,
		WaitForHosts:   waitForHosts,
		Axes:           axes,
		Behaviors:      behaviors,
		JSONReportPath: jsonReportPath,
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

// TODO: return an error if input is malformed.
func parseBehavior(d string) Behavior {
	pair := strings.SplitN(d, "=", 2)
	key := strings.ToLower(pair[0])

	values := strings.Split(pair[1], ",")
	values = trimCollection(values)
	clientAxis := values[0]
	values = values[1:]
	sort.Strings(values)

	behavior := Behavior{
		Name:       key,
		ClientAxis: clientAxis,
		ParamsAxes: values,
	}
	return behavior
}

func parseSkipBehavior(d string) (string, []Filter, error) {
	pair := strings.SplitN(d, "=", 2)
	if len(pair) != 2 {
		return "", nil, fmt.Errorf("missing '=' in the input: %q", d)
	}
	behaviorName := strings.ToLower(pair[0])
	rawFilters := strings.Split(pair[1], ",")
	filters := make([]Filter, 0, len(rawFilters))
	for _, rawFilter := range rawFilters {
		rawMatches := strings.Split(rawFilter, "+")
		filter := Filter{
			AxisMatches: make([]AxisMatch, 0, len(rawMatches)),
		}
		for _, rawMatch := range rawMatches {
			tuple := strings.SplitN(rawMatch, ":", 2)
			tuple = trimCollection(tuple)
			if len(tuple) != 2 {
				return "", nil, fmt.Errorf("match %q in input %s is not of form 'key:value'", rawMatch, d)
			}
			axisName := tuple[0]
			axisValue := tuple[1]
			filter.AxisMatches = append(filter.AxisMatches, AxisMatch{Name: axisName, Value: axisValue})
		}
		filters = append(filters, filter)
	}
	return behaviorName, filters, nil
}

func parseAxis(d string) Axis {
	pair := strings.SplitN(d, "=", 2)
	key := strings.ToLower(pair[0])
	values := strings.Split(pair[1], ",")
	values = trimCollection(values)
	sort.Strings(values)

	axis := Axis{
		Name:   key,
		Values: values,
	}

	return axis
}

// TODO: Return an error if input is malformed.
func validateConfig(config *Config) error {
	axes := config.Axes.Index()
	for _, behavior := range config.Behaviors {
		if _, ok := axes[behavior.ClientAxis]; !ok {
			return errors.New("can't find AXIS environment for: " + behavior.ClientAxis)
		}
		for _, param := range behavior.ParamsAxes {
			if _, ok := axes[param]; !ok {
				return errors.New("can't find AXIS environment for: " + param)
			}
		}
	}
	return nil
}

func trimCollection(in []string) []string {
	ret := make([]string, 0, len(in))
	for _, v := range in {
		if v == "" {
			continue
		}
		ret = append(ret, strings.Trim(v, " "))
	}
	return ret
}
