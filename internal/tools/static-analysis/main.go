// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// This tool runs project-specific static analysis rules that cannot be implemented
// using standard linters like golangci-lint. These rules require domain knowledge of
// the provider's internals, such as inspecting typed SDK resource models via runtime
// reflection or enforcing code style conventions specific to this project.
//
// Usage:
//
//	go run internal/tools/static-analysis/main.go                           # run all rules
//	go run internal/tools/static-analysis/main.go -rules=combinedIfErr     # run a specific rule
//	go run internal/tools/static-analysis/main.go -fail-on-error=false     # log errors without failing

package main

import (
	"flag"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/static-analysis/rules"
)

var allRules = map[string]rules.Rule{
	rules.TypedSDKBitCheck{}.Name():   rules.TypedSDKBitCheck{},
	rules.CombinedIfErrCheck{}.Name(): rules.CombinedIfErrCheck{},
}

func main() {
	f := flag.NewFlagSet("staticAnalysis", flag.ExitOnError)

	rulesToCheck := f.String("rules", "all", "Comma separated list of rules to run. Defaults to all. ")
	failOnError := f.Bool("fail-on-error", true, "If set to true will fail on error, otherwise will only log. Defaults to true.")

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	if len(*rulesToCheck) == 0 {
		log.Fatalf("no rules specified")
	}
	specifiedRules := strings.Split(*rulesToCheck, ",")

	// If `all` is in the list, just reset it to `all`
	if slices.Contains(specifiedRules, "all") {
		specifiedRules = []string{"all"}
	}

	errors := make([]error, 0)
	for _, rule := range specifiedRules {
		if strings.EqualFold(rule, "all") {
			for _, r := range allRules {
				errors = append(errors, r.Run()...)
			}
		}

		if r, ok := allRules[rule]; ok {
			errors = append(errors, r.Run()...)
		}
	}

	if len(errors) > 0 {
		log.Printf("Static analysis found %d error(s):\n", len(errors))
		for _, err := range errors {
			log.Printf("  - %s\n", err)
		}
		if *failOnError {
			os.Exit(1)
		}
	}
}
