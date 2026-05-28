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
	case "setup verify":
		printSetupVerify()
	default:
		usage()
		os.Exit(1)
	}
}

func printSetupPlan() {
	plan := setup.DefaultInstallPlan()
	printJSON(plan, "failed to encode install plan")
}

func printSetupVerify() {
	report := setup.VerifyLocalKubernetesOperator()
	printJSON(report, "failed to encode setup verification report")

	if !report.OK {
		os.Exit(2)
	}
}

func printJSON(value any, failureMessage string) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", failureMessage, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  headlamp setup plan")
	fmt.Println("  headlamp setup verify")
}
