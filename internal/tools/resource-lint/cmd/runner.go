package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

// ExitCode represents program exit codes
type ExitCode int

const (
	ExitSuccess     ExitCode = 0 // No issues found
	ExitIssuesFound ExitCode = 1 // Lint issues found
	ExitError       ExitCode = 2 // Tool error
)

type Runner struct {
	Config *Config
}

// NewRunner creates a new Runner with the given config
func NewRunner(cfg *Config) *Runner {
	return &Runner{
		Config: cfg,
	}
}

// Run executes the linter and returns an exit code
func (r *Runner) Run(ctx context.Context) ExitCode {
	loaderOpts := loader.LoaderOptions{
		NoFilter:   r.Config.NoFilter,
		RemoteName: r.Config.RemoteName,
		BaseBranch: r.Config.BaseBranch,
		DiffFile:   r.Config.DiffFile,
	}

	_, err := loader.LoadChanges(loaderOpts)
	if err != nil {
		log.Printf("Warning: failed to load changed lines filter: %v", err)
	}

	// Determine package patterns to analyze
	patterns := r.Config.Patterns
	if loader.IsEnabled() {
		files, lines := loader.GetStats()
		log.Printf("Changed lines filter: tracking %d files with %d changed lines", files, lines)

		// If change tracking is enabled and no patterns specified, use changed packages
		if len(r.Config.Patterns) == 0 {
			changedPackages := loader.GetChangedPackages()
			if len(changedPackages) > 0 {
				patterns = changedPackages
				log.Printf("Auto-detected %d changed packages:", len(patterns))
				for _, pkg := range patterns {
					log.Printf("  %s", pkg)
				}
			}
		}
	}

	// Validate we have patterns to analyze
	if len(patterns) == 0 {
		log.Println("Error: no packages to analyze")
		return ExitError
	}

	log.Printf("Loading packages...")
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Printf("Error: failed to load packages: %v", err)
		return ExitError
	}

	// Check for package loading errors
	var hasLoadErrors bool
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			log.Printf("Error: failed to load package: %v", err)
			hasLoadErrors = true
		}
	})
	if hasLoadErrors {
		return ExitError
	}

	// Provide loaded packages to analyzers for cross-package schema resolution
	helper.SetGlobalPackages(pkgs)

	log.Printf("Running analysis...")
	graph, err := checker.Analyze(passes.AllChecks, pkgs, nil)
	if err != nil {
		log.Printf("Error: analysis failed: %v", err)
		return ExitError
	}

	// Report diagnostics
	foundIssues := r.reportDiagnostics(graph)
	if foundIssues {
		return ExitIssuesFound
	}

	log.Printf("✓ Analysis completed successfully with no issues found")
	return ExitSuccess
}

// reportDiagnostics reports all diagnostics and returns true if any issues were found
func (r *Runner) reportDiagnostics(graph *checker.Graph) bool {
	var foundIssues bool
	var issueCount int

	for act := range graph.All() {
		if act.Err != nil {
			fmt.Printf("%s: %v\n", act.Package.PkgPath, act.Err)
			foundIssues = true
			continue
		}

		for _, diag := range act.Diagnostics {
			foundIssues = true
			issueCount++
			fmt.Printf("%s: %s\n", act.Package.Fset.Position(diag.Pos), diag.Message)
		}
	}

	if foundIssues {
		fmt.Printf("Found %d issue(s)\n", issueCount)
	}

	return foundIssues
}
