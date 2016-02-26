package plan

import (
	"os"
	"strings"
)

const clientsKey = "CROSSDOCK_CLIENTS"
const axisKeyPrefix = "CROSSDOCK_AXIS_"

// ReadFromEnviron creates a Plan by looking for CROSSDOCK_ environment vars
func ReadFromEnviron() Plan {
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

	plan := Plan{
		Clients: clients,
		Axes:    axes,
	}

	return plan
}
