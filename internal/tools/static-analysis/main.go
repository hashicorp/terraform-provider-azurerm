// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	rules.TypedSDKBitCheck{}.Name(): rules.TypedSDKBitCheck{},
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
		if *failOnError {
			log.Fatalf("failed to run rules: %v", errors)
		} else {
			log.Printf("failed to run rules: %v", errors)
			os.Exit(0)
		}
	}
}
