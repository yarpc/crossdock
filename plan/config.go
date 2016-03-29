package plan

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultCallTimeout = 5

// ReadConfigFromEnviron creates a Config by looking for CROSSDOCK_ environment vars
func ReadConfigFromEnviron() *Config {
	const (
		callTimeoutKey    = "CROSSDOCK_CALL_TIMEOUT"
		waitKey           = "CROSSDOCK_WAIT_FOR"
		axisKeyPrefix     = "CROSSDOCK_AXIS_"
		behaviorKeyPrefix = "CROSSDOCK_BEHAVIOR_"
	)
	callTimeout, _ := strconv.Atoi(os.Getenv(callTimeoutKey))
	if callTimeout == 0 {
		callTimeout = defaultCallTimeout
	}

	WaitForHosts := trimCollection(strings.Split(os.Getenv(waitKey), ","))
	axes := make(map[string]Axis)
	behaviors := make(map[string]Behavior)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, axisKeyPrefix) {
			axis := parseAxis(strings.TrimPrefix(e, axisKeyPrefix))
			axes[axis.Name] = axis
		} else if strings.HasPrefix(e, behaviorKeyPrefix) {
			behavior := parseBehavior(strings.TrimPrefix(e, behaviorKeyPrefix))
			behaviors[behavior.Name] = behavior
		}
	}
	config := &Config{
		CallTimeout:  time.Duration(callTimeout),
		WaitForHosts: WaitForHosts,
		Axes:         axes,
		Behaviors:    behaviors,
	}

	validateConfig(config)
	return config
}

func parseBehavior(d string) Behavior {
	pair := strings.SplitN(d, "=", 2)
	key := strings.ToLower(pair[0])
	values := strings.Split(pair[1], ",")
	values = trimCollection(values)
	behavior := Behavior{
		Name:    key,
		Clients: values[0],
		Params:  values[1:],
	}

	return behavior
}

func parseAxis(d string) Axis {
	pair := strings.SplitN(d, "=", 2)
	key := strings.ToLower(pair[0])
	values := strings.Split(pair[1], ",")
	values = trimCollection(values)

	axis := Axis{
		Name:   key,
		Values: values,
	}

	return axis
}

func validateConfig(config *Config) {
	for _, behavior := range config.Behaviors {
		if _, ok := config.Axes[behavior.Clients]; !ok {
			log.Fatal("Can't find AXIS environment for: " + behavior.Clients)
		}
		for _, param := range behavior.Params {
			if _, ok := config.Axes[param]; !ok {
				log.Fatal("Can't find AXIS environment for: " + param)
			}
		}
	}
}

func trimCollection(in []string) []string {
	ret := make([]string, len(in))
	for i, v := range in {
		ret[i] = strings.Trim(v, " ")
	}
	return ret
}
