package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/cmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Parse configuration from flags
	cfg, err := cmd.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: resource-lint [flags] <package patterns>\n")
		return 3
	}

	// Handle help flag
	if cfg.ShowHelp {
		cfg.PrintHelp()
		return 0
	}

	// Handle list checks flag
	if cfg.ListChecks {
		cmd.PrintChecks()
		return 0
	}

	// Create and run the linter
	runner := cmd.NewRunner(cfg)
	exitCode := runner.Run(context.Background())

	return int(exitCode)
}
