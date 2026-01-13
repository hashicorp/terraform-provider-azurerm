package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes"
)

// Config holds all configuration options for the linter
type Config struct {
	// Command options
	Patterns   []string
	ShowHelp   bool
	ListChecks bool

	// Loader options
	NoFilter   bool
	RemoteName string
	BaseBranch string
	DiffFile   string

	// Internal: flagSet for help printing
	flagSet *flag.FlagSet
}

// ParseFlags parses command line flags and returns a Config
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet("resource-lint", flag.ExitOnError)

	// Config struct to populate
	cfg := &Config{flagSet: fs}

	// Command flags
	fs.BoolVar(&cfg.ShowHelp, "help", false, "show help message")
	fs.BoolVar(&cfg.ListChecks, "list", false, "list all available checks")

	// Loader flags
	fs.BoolVar(&cfg.NoFilter, "no-filter", false, "disable change filtering, analyze all files")
	fs.StringVar(&cfg.RemoteName, "remote", "", "git remote name (auto-detect: origin > upstream)")
	fs.StringVar(&cfg.BaseBranch, "base", "", "base branch (auto-detect from git config or 'main')")
	fs.StringVar(&cfg.DiffFile, "diff", "", "read diff from file instead of git")

	fs.Usage = func() {
		cfg.PrintHelp()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg.Patterns = fs.Args()

	return cfg, nil
}

// PrintHelp prints the help message
func (c *Config) PrintHelp() {
	fmt.Println(`resource-lint - AzureRM Provider resource linting tool

Usage:
  go run ./internal/tools/resource-lint [flags] <package patterns>

Examples:
  go run ./internal/tools/resource-lint ./internal/services/compute/...
  go run ./internal/tools/resource-lint --pr=12345
  go run ./internal/tools/resource-lint --diff=changes.txt
  go run ./internal/tools/resource-lint --no-filter ./internal/services/...

Flags:`)
	c.flagSet.PrintDefaults()
}

// PrintChecks prints all available checks
func PrintChecks() {
	fmt.Println("Available checks:")
	for _, analyzer := range passes.AllChecks {
		title := strings.Split(analyzer.Doc, "\n")[0]
		fmt.Printf("  %-10s  %s\n", analyzer.Name, title)
	}
}
