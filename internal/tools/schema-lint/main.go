// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Command schema-lint is a modular, pluggable linter for the AzureRM provider's
// resource and data source schemas.
//
// In the spirit of a markdown linter, each check is an independent rule that can
// be enabled, disabled or re-configured via a .schema-lint.json config file or
// CLI flags. Run `schema-lint list` to see the available rules.
//
// Usage:
//
//	go run ./internal/tools/schema-lint check                       # lint with defaults
//	go run ./internal/tools/schema-lint check -rules=SL001,SL005    # run only specific rules
//	go run ./internal/tools/schema-lint check -disable=SL001        # disable a rule
//	go run ./internal/tools/schema-lint check -fix                  # include suggested fixes
//	go run ./internal/tools/schema-lint check -diff base.json       # only lint properties added since base.json
//	go run ./internal/tools/schema-lint check -format=json          # machine readable output
//	go run ./internal/tools/schema-lint list                        # list all rules
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/engine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/report"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-lint/rules"
)

func main() {
	args := os.Args[1:]

	cmd := "check"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "list":
		listRules()
	case "check":
		os.Exit(runCheck(args))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		printHelp()
		os.Exit(2)
	}
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("schema-lint check", flag.ExitOnError)
	configPath := fs.String("config", "", "path to a .schema-lint.json config file (defaults to .schema-lint.json when present)")
	rulesFlag := fs.String("rules", "all", "comma-separated rule IDs to run, or 'all'")
	disableFlag := fs.String("disable", "", "comma-separated rule IDs to disable")
	resourceFlag := fs.String("resource", "", "comma-separated resource/data source type names to include")
	skipResourceFlag := fs.String("skip-resource", "", "comma-separated resource/data source type names to skip")
	formatFlag := fs.String("format", "text", "output format: text or json")
	failOnError := fs.Bool("fail-on-error", true, "exit non-zero when any error-severity finding is present")
	fixFlag := fs.Bool("fix", false, "include suggested fixes for fixable findings")
	diffFlag := fs.String("diff", "", "path to a base schema dump (from schema-api -export); only report findings on properties added since the base")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		return 2
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		return 2
	}

	// In diff mode, load the base schema so that only newly-added properties are
	// held to the rules.
	var baseSchema *providerschema.ProviderSchemaJSON
	if *diffFlag != "" {
		wrapper, err := providerschema.LoadWrapperFromFile(*diffFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading base schema %q: %v\n", *diffFlag, err)
			return 2
		}
		baseSchema = wrapper.ProviderSchema
	}

	// CLI resource filters override the config file when provided.
	if v := splitList(*resourceFlag); len(v) > 0 {
		cfg.IncludeResources = v
	}
	if v := splitList(*skipResourceFlag); len(v) > 0 {
		cfg.SkipResources = v
	}

	only := splitList(*rulesFlag)
	if len(only) == 1 && strings.EqualFold(only[0], "all") {
		only = nil
	}

	linter := engine.New(nil, engine.Options{
		Config:       cfg,
		OnlyRules:    only,
		DisableRules: splitList(*disableFlag),
		SuggestFixes: *fixFlag,
		BaseSchema:   baseSchema,
	})

	findings, err := linter.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running linter: %v\n", err)
		return 2
	}

	if err := report.Write(os.Stdout, findings, report.Format(strings.ToLower(*formatFlag))); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		return 2
	}

	if *failOnError && report.Summarize(findings).Errors > 0 {
		return 1
	}

	return 0
}

func listRules() {
	fmt.Println("Available schema-lint rules:")
	for _, r := range rules.AllRules {
		fixable := ""
		if rules.Fixable(r) {
			fixable = " [fixable]"
		}
		fmt.Printf("  %-6s [%-7s] %s%s\n           %s\n", r.ID(), r.DefaultSeverity(), r.Name(), fixable, r.Description())
	}
}

func printHelp() {
	fmt.Print(`USAGE: schema-lint [COMMAND] [OPTIONS]

COMMANDS:
  check   lint the provider schema (default)
  list    show the available rules

Run "schema-lint check -h" for the available options.
`)
}

func splitList(in string) []string {
	if strings.TrimSpace(in) == "" {
		return nil
	}

	parts := strings.Split(in, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}

	return out
}
