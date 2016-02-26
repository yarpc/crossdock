package plan

import (
	"os"
	"strings"
)

// ReadFromEnviron creates a Plan based on environment vars
func ReadFromEnviron() Plan {
	config := ReadConfigFromEnviron()
	plan := New(config)
	return plan
}

// ReadConfigFromEnviron creates a Config by looking for CROSSDOCK_ environment vars
func ReadConfigFromEnviron() Config {
	const clientsKey = "CROSSDOCK_CLIENTS"
	const axisKeyPrefix = "CROSSDOCK_AXIS_"

	clients := strings.Split(os.Getenv(clientsKey), ",")
	var axes []Axis

	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, axisKeyPrefix) {
			continue
		}
		d := strings.TrimPrefix(e, axisKeyPrefix)

		pair := strings.SplitN(d, "=", 2)
		key := strings.ToLower(pair[0])
		value := strings.Split(pair[1], ",")

		axis := Axis{
			Name:   key,
			Values: value,
		}
		axes = append(axes, axis)
	}

	config := Config{
		Clients: clients,
		Axes:    axes,
	}

	return config
}
