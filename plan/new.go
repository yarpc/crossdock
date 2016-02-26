package plan

// New creates a Plan given a Config
func New(config Config) Plan {
	return Plan{
		Config:    config,
		TestCases: collect(config),
	}
}

func collect(config Config) []TestCase {
	var cases []TestCase
	for _, client := range config.Clients {
		clientTestCases := createCasesForClient(client, config.Axes, make(map[string]string))
		cases = append(cases, clientTestCases...)
	}
	return cases
}

func createCasesForClient(client string, axes []Axis, args Arguments) []TestCase {
	if len(axes) == 0 {
		return []TestCase{{
			Client:    client,
			Arguments: copyArgs(args),
		}}
	}
	var testCases []TestCase
	axis := axes[0]
	for _, p := range axis.Values {
		args[axis.Name] = p
		testCases = append(testCases, createCasesForClient(client, axes[1:], args)...)
	}
	return testCases
}

func copyArgs(args Arguments) Arguments {
	copy := make(map[string]string)
	for k, v := range args {
		copy[k] = v
	}
	return copy
}
