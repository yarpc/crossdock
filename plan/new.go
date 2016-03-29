package plan

// New creates a Plan given a Config
func New(config *Config) *Plan {
	plan := &Plan{
		Config: config,
	}
	plan.TestCases = buildTestCases(plan)
	return plan
}

func buildTestCases(plan *Plan) []TestCase {
	var testCases []TestCase
	for _, behavior := range plan.Config.Behaviors {
		for _, client := range plan.Config.Axes[behavior.Clients].Values {
			combos := recurseCombinations(0, plan, client, behavior, map[string]string{"behavior": behavior.Name})
			testCases = append(testCases, combos...)
		}
	}
	return testCases
}

func recurseCombinations(level int, plan *Plan, client string, behavior Behavior, args Arguments) []TestCase {
	if level == len(behavior.Params) {
		return []TestCase{{
			Plan:      plan,
			Client:    client,
			Arguments: copyArgs(args),
		}}
	}
	var testCases []TestCase
	param := behavior.Params[level]
	for _, axis := range plan.Config.Axes[param].Values {
		args[param] = axis
		testCases = append(testCases, recurseCombinations(level+1, plan, client, behavior, args)...)
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
