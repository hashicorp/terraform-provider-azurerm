// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/check"
)

func printHelp() {
	text := `USAGE: go run main.go [CMD] [OPTIONS]
CMD:
  check:	check documents and print the error information
  fix:	 	check and try to fix existing errors

OPTIONS:
`
	fmt.Printf("%s\n", text)
}

var (
	cmd          string
	dryRun       = true
	resource     string
	service      string
	skipResource string
	skipService  string
)

func parseArgs() {
	fs := flag.NewFlagSet("azdoc-check", flag.ExitOnError)
	fs.StringVar(&resource, "resource", os.Getenv("ONLY_RESOURCE"), "a list of resource names to check")
	fs.StringVar(&skipResource, "skip-resource", os.Getenv("SKIP_RESOURCE"), "a list of resource names to skip the check")
	fs.StringVar(&service, "service", os.Getenv("ONLY_SERVICE"), "a list of services names to check")
	fs.StringVar(&skipService, "skip-service", os.Getenv("SKIP_SERVICE"), "a list of service names to skip the check")

	fs.Usage = func() {
		printHelp()
		fs.PrintDefaults()
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		cmd = os.Args[1]
		switch cmd {
		case "check":
			_ = fs.Parse(os.Args[2:])
		case "fix":
			dryRun = false
			_ = fs.Parse(os.Args[2:])
		default:
			fs.Usage()
		}
	}
}

func main() {
	parseArgs()

	result := check.DiffAll(check.AzurermAllResources(service, skipService, resource, skipResource), dryRun)
	if !result.HasDiff() {
		log.Printf("document linter runs success, time costs: %v", result.CostTime())
		return
	}

	log.Printf("%s\n", result.ToString())

	if cmd == "fix" {
		if err := result.FixDocuments(); err != nil {
			log.Fatalf("error occurs when trying to fix documents: %v", err)
		}
	}
	os.Exit(1)
}
