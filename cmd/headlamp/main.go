package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/unboxd-cloud/headlamp-agent/pkg/setup"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	section := os.Args[1]
	command := os.Args[2]

	switch section + " " + command {
	case "setup plan":
		printSetupPlan()
	default:
		usage()
		os.Exit(1)
	}
}

func printSetupPlan() {
	plan := setup.DefaultInstallPlan()
	encoded, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode install plan: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  headlamp setup plan")
}
