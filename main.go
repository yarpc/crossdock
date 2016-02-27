package main

import (
	"fmt"
	"time"

	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/plan"
)

func main() {
	fmt.Printf("\nCrossdock starting...\n\n")
	config := plan.ReadConfigFromEnviron()
	plan := plan.New(config)

	fmt.Printf("Waiting on CROSSDOCK_CLIENTS=%v\n\n", config.Clients)
	execute.Wait(config.Clients, time.Duration(30)*time.Second)

	fmt.Printf("\nExecuting Matrix...\n\n")
	results := execute.Run(plan)
	Output(results)
}
