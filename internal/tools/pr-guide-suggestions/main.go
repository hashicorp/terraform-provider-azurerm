// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Command pr-guide-suggestions is a fast, AST-based linter for the AzureRM
// provider's resource and data source schemas. It extends the tfproviderlint
// helper libraries and analysis.Analyzer model, parsing Go source directly (no
// provider compilation or JSON schema rendering) so it can lint a single file or
// only the schema properties a pull request changes.
//
// Usage:
//
//	pr-guide-suggestions list
//	pr-guide-suggestions check ./internal/services/foo/foo_resource.go
//	pr-guide-suggestions check ./internal/services/foo
//	pr-guide-suggestions check -diff-base origin/main
//	pr-guide-suggestions check -format json ./internal/services/foo
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/lint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/rules"
)

// defaultPathspecs restrict diff mode to resource and data source files,
// matching the pr-guide-suggestions workflow's trigger paths.
var defaultPathspecs = []string{
	":(glob)internal/services/**/*_resource.go",
	":(glob)internal/services/**/*_data_source.go",
	":(glob)internal/services/**/*_datasource.go",
}

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
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	root := fs.String("C", ".", "repository root to resolve targets and run git from")
	format := fs.String("format", "text", "output format: text or json")
	fix := fs.Bool("fix", false, "apply auto-fixable fixes (property renames) to files in place")
	diffBase := fs.String("diff-base", "", "only report findings on lines added since this git ref (e.g. origin/main)")
	rulesFlag := fs.String("rules", "", "comma-separated rule IDs to run (default all)")
	disableFlag := fs.String("disable", "", "comma-separated rule IDs to disable")
	failOnError := fs.Bool("fail-on-error", true, "exit non-zero when any error-severity finding is present")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving root %q: %v\n", *root, err)
		return 2
	}

	opts := lint.Options{
		Only:    splitSet(*rulesFlag),
		Disable: splitSet(*disableFlag),
	}

	targets := fs.Args()
	if *diffBase != "" {
		changes, err := lint.Diff(absRoot, *diffBase, defaultPathspecs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error computing diff against %q: %v\n", *diffBase, err)
			return 2
		}
		opts.Changes = changes
		targets = changes.ChangedFiles()
		if len(targets) == 0 {
			fmt.Fprintln(os.Stderr, "no changed resource or data source schema files; nothing to lint")
			return 0
		}
	}

	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "error: no targets; provide files/directories or use -diff-base")
		printHelp()
		return 2
	}

	findings, err := lint.Run(targets, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running linter: %v\n", err)
		return 2
	}

	if *fix {
		applied, err := lint.Apply(findings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error applying fixes: %v\n", err)
			return 2
		}
		lint.WriteApplied(os.Stderr, applied, absRoot)
		// The renames are now applied; report only what still needs manual work.
		findings = lint.Unfixed(findings)
	}

	switch strings.ToLower(*format) {
	case "json":
		if err := lint.WriteJSON(os.Stdout, findings); err != nil {
			fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
			return 2
		}
	default:
		lint.WriteText(os.Stdout, findings, absRoot)
	}

	if *failOnError && lint.HasErrors(findings) {
		return 1
	}
	return 0
}

func listRules() {
	fmt.Println("ID     Severity  Fixable  Name")
	for _, r := range rules.Registry {
		fixable := ""
		if r.Fixable {
			fixable = "yes"
		}
		fmt.Printf("%-6s %-9s %-8s %s\n", r.ID, r.Severity, fixable, r.Name)
	}
}

func printHelp() {
	fmt.Fprintln(os.Stderr, `pr-guide-suggestions - AST-based AzureRM schema linter

Commands:
  check [flags] [files/dirs...]   lint the given targets (default command)
  list                            show the available rules

Check flags:
  -C <dir>            repository root (default ".")
  -format text|json   output format (default "text")
  -fix                apply auto-fixable fixes (property renames) in place
  -diff-base <ref>    lint only lines added since <ref> (e.g. origin/main)
  -rules <ids>        comma-separated rule IDs to run
  -disable <ids>      comma-separated rule IDs to disable
  -fail-on-error      exit non-zero on error findings (default true)`)
}

func splitSet(s string) map[string]bool {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	out := map[string]bool{}
	for _, part := range strings.Split(s, ",") {
		if p := strings.ToUpper(strings.TrimSpace(part)); p != "" {
			out[p] = true
		}
	}
	return out
}
