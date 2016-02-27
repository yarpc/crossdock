package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yarpc/crossdock/execute"
	"github.com/yarpc/crossdock/output"
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

	if err := output.Stream(results); err != nil {
		fmt.Println("Tests did not pass!")
		os.Exit(1)
	}
}
