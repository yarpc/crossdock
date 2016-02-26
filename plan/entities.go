package plan

// Axis represents combinational args to be passed to the test clients
type Axis struct {
	Name   string
	Values []string
}

// Plan describes the entirety of the test program
type Plan struct {
	Clients []string
	Axes    []Axis
}
