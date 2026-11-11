package main

import (
	"fmt"
	"os"

	"github.com/ogormans-deptstack/orca-temporal-bridge/cmd/orca-temporal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
