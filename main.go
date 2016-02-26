package main

import (
	"fmt"
	"time"

	"github.com/yarpc/crossdock/plan"
)

func main() {
	fmt.Printf("\nCrossdock starting...\n\n")

	plan := plan.ReadFromEnviron()

	fmt.Printf("Waiting on CROSSDOCK_CLIENTS=%v\n\n", plan.Clients)
	Wait(plan.Clients, time.Duration(30)*time.Second)

	fmt.Printf("\nExecuting Matrix...\n\n")
	results := Execute(plan)

	Output(results)
}
