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
	plan := plan.New(plan.ReadConfigFromEnviron())

	fmt.Printf("Waiting on CROSSDOCK_WAIT_FOR=%v\n\n", plan.Config.WaitForHosts)
	execute.Wait(plan.Config.WaitForHosts, time.Duration(30)*time.Second)

	fmt.Printf("\nExecuting Matrix...\n\n")
	results := execute.Run(plan)

	if err := output.Stream(results); err != nil {
		fmt.Printf("\nTests did not pass!\n\n")
		os.Exit(1)
	}
	fmt.Printf("\nTests passed!\n\n")
}
