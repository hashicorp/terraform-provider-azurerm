package cmd

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	"github.com/spf13/cobra"
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
	Rules []string
}

// FlagsScaffold should contain the flags required by `website-scaffold` once that functionality
// has been moved into this tool.
type FlagsScaffold struct{}

var flags Flags

func configureFlags(root *cobra.Command) {
	root.PersistentFlags().BoolVarP(&flags.Debug, "debug", "d", false, "output debug logs, defaults to false")
	root.PersistentFlags().StringVarP(&flags.ProviderDirectory, "provider-directory", "p", "", "provider directory path, can be omitted if provider directory is the current working directory")
	root.PersistentFlags().StringVarP(&flags.Service, "service", "s", "", "service to filter the operation to")
	root.PersistentFlags().StringVarP(&flags.Resource, "resource", "r", "", "resource to filter the operation to")

	for _, c := range root.Commands() {
		switch c.Name() {
		case "validate", "fix":
			c.PersistentFlags().StringSliceVar(&flags.Linter.Rules, "rules", util.MapKeys2Slice(rule.Registration), "A comma separated list of rule IDs, if not specified, all rules will be run")
		}
	}
}
