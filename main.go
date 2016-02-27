package main

import (
	"fmt"
	"time"

	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
)

func main() {
	fmt.Printf("\nCrossdock starting...\n\n")
	plan := plan.New(plan.ReadConfigFromEnviron())

	fmt.Printf("Waiting on CROSSDOCK_CLIENTS=%v\n\n", plan.Config.Clients)
	Wait(plan.Config.Clients, time.Duration(30)*time.Second)

	fmt.Printf("\nExecuting Matrix...\n\n")
	results := execute.Run(plan)
	Output(results)
}
