package plan

// Config describes the unstructured test plan
type Config struct {
	Clients []string
	Axes    []Axis
}

// Axis represents combinational args to be passed to the test clients
type Axis struct {
	Name   string
	Values []string
}

// Plan describes the entirety of the test program
type Plan struct {
	TestCases []TestCase
}

// TestCase represents the request made to test clients.
type TestCase struct {
	Client    string
	Arguments Arguments
}

// Arguments represents custom args to pass to test client.
type Arguments map[string]string
