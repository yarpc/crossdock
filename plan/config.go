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
		if strings.HasPrefix(e, axisKeyPrefix) {
			axis := parseAxis(strings.TrimPrefix(e, axisKeyPrefix))
			axes = append(axes, axis)
		} else if strings.HasPrefix(e, skipKeyPrefix) {
			key, filter := parseSkipBehavior(strings.TrimPrefix(e, skipKeyPrefix))
			filterMap[key] = filter
		} else if strings.HasPrefix(e, behaviorKeyPrefix) {
			behavior := parseBehavior(strings.TrimPrefix(e, behaviorKeyPrefix))
			behaviors = append(behaviors, behavior)
		}
	}
	sort.Sort(axes)
	sort.Sort(behaviors)
	err := behaviors.validateAndApplyFilters(filterMap)
	if err != nil {
		return nil, err
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

// Skip string is of format SKIP_RUN=axis1:value1+axis2:value2,axis1:value3 where
// SKIP_ is followed by behavior name and multiple filters are separated by comma.
// Each filter is a logical AND and complete match is supported.
func parseSkipBehavior(d string) (string, []Filter) {
	pair := strings.SplitN(d, "=", 2)
	key := strings.ToLower(pair[0])
	rawFilters := strings.Split(pair[1], ",")
	var filters []Filter
	for _, rawFilter := range rawFilters {
		strMatches := strings.Split(rawFilter, "+")
		filter := Filter{
			AxisMatches: map[string]string{},
		}
		for _, strMatch := range strMatches {
			tuple := strings.Split(strMatch, ":")
			tuple = trimCollection(tuple)
			filter.AxisMatches[tuple[0]] = tuple[1]
		}
		filters = append(filters, filter)
	}
	return key, filters
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

func validateConfig(config *Config) error {
	axes := config.Axes.Index()
	for _, behavior := range config.Behaviors {
		if _, ok := axes[behavior.ClientAxis]; !ok {
			return errors.New("Can't find AXIS environment for: " + behavior.ClientAxis)
		}
		for _, param := range behavior.ParamsAxes {
			if _, ok := axes[param]; !ok {
				return errors.New("Can't find AXIS environment for: " + param)
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
