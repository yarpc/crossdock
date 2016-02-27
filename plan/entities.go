package plan

// Axis represents combinational args to be passed to the test clients
type Axis struct {
	Name   string
	Values []string
}

// Config describes the unstructured test plan
type Config struct {
	Clients []string
	Axes    []Axis
}

// Arguments represents custom args to pass to test client.
type Arguments map[string]string

// TestCase represents the request made to test clients.
type TestCase struct {
	Client    string
	Arguments Arguments
}

// Plan describes the entirety of the test program
type Plan struct {
	TestCases []TestCase
}
