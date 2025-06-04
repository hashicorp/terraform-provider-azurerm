package cmd

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/cmd/values"
	"github.com/spf13/pflag"
)

type Flags struct {
	Debug             bool
	ProviderDirectory string
	Service           string
	Resource          string

	Linter   FlagsLinter
	Scaffold FlagsScaffold
}

type FlagsLinter struct {
	Rules string
}

// FlagsScaffold should contain the flags required by `website-scaffold` once that functionality
// has been moved into this tool.
type FlagsScaffold struct{}

var flags Flags

func configureFlags() {
	rootCmd.PersistentFlags().BoolVarP(&flags.Debug, "debug", "d", false, "output debug logs, defaults to false")
	rootCmd.PersistentFlags().StringVarP(&flags.ProviderDirectory, "provider-directory", "p", "", "provider directory path, can be omitted if provider directory is the current working directory")
	rootCmd.PersistentFlags().StringVarP(&flags.Service, "service", "s", "", "service to filter the operation to")
	rootCmd.PersistentFlags().StringVarP(&flags.Resource, "resource", "r", "", "resource to filter the operation to")

	rulesValue := values.NewStringValue("", &flags.Linter.Rules)
	rulesFlag := pflag.Flag{
		Name:     "rules",
		Usage:    "A comma separated list of rule IDs, if not specified, all rules will be run",
		Value:    rulesValue,
		DefValue: rulesValue.String(),
	}

	validateCmd.Flags().AddFlag(&rulesFlag)
	fixCmd.Flags().AddFlag(&rulesFlag)
}
