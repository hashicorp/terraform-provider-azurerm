package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func Make() *cobra.Command {
	fs := afero.NewOsFs()

	rootCmd := &cobra.Command{
		Use:   "documentfmt [validate|fix|scaffold]",
		Short: "A small tool to validate provider documentation.",
		Long:  `A small tool to validate provider documentation based on a set of custom rules. It can also fix most found issues.`,
	}

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates documentation",
		PreRun: func(cmd *cobra.Command, args []string) {
			util.InitLogger(flags.Debug)
			validateProviderDirectoryAccess(fs)
		},
		Run: func(cmd *cobra.Command, args []string) {
			resources := data.GetAllTerraformNodeData(fs, flags.ProviderDirectory, flags.Service, flags.Resource)

			v := validator.Validator{}
			for _, r := range resources {
				v.Run(flags.Linter.Rules, r, false)
			}

			// separate from loop above to prevent debug messages from mixing in
			errCount, resourceWithErrCount := 0, 0
			for _, r := range resources {
				if l := len(r.Errors); l > 0 {
					resourceWithErrCount++
					errCount += l
					printErrors(r)
				}
			}

			if errCount > 0 {
				errStr, resourceStr := "error", "resource"
				if errCount > 1 {
					errStr += "s"
				}

				if resourceWithErrCount > 1 {
					resourceStr += "s"
				}

				fmt.Printf(util.Red("Found %d %s in %d %s\n"), errCount, errStr, resourceWithErrCount, resourceStr)
				os.Exit(1)
			}

			fmt.Print(util.GreenBold("Found no errors\n"))
		},
	}

	fixCmd := &cobra.Command{
		Use:   "fix",
		Short: "Attempts to fix documentation",
		PreRun: func(cmd *cobra.Command, args []string) {
			util.InitLogger(flags.Debug)
			validateProviderDirectoryAccess(fs)
		},
		Run: func(cmd *cobra.Command, args []string) {
			resources := data.GetAllTerraformNodeData(fs, flags.ProviderDirectory, flags.Service, flags.Resource)

			v := validator.Validator{}
			errCount, resourceWithErrCount := 0, 0
			for _, r := range resources {
				v.Run(flags.Linter.Rules, r, true)

				if l := len(r.Errors); l > 0 {
					resourceWithErrCount++
					errCount += l
				}

				if r.Document.HasChange {
					if err := r.Document.Write(fs); err != nil {
						if err != nil {
							log.WithFields(log.Fields{
								"resource": r.Name,
								"path":     r.Document.Path,
								"error":    err,
							}).Error("writing changes to the documentation file")
						}
					}
				}
			}

			if errCount > 0 {
				errStr, resourceStr := "error", "resource"
				if errCount > 1 {
					errStr += "s"
				}

				if resourceWithErrCount > 1 {
					resourceStr += "s"
				}
				fmt.Printf(util.Red("Found %d %s in %d %s and applied fixes where possible, please review the changes\n"), errCount, errStr, resourceWithErrCount, resourceStr)
				os.Exit(1)
			} else {
				fmt.Print(util.GreenBold("Found no errors\n"))
			}
		},
	}

	scaffoldCmd := &cobra.Command{
		Use:   "scaffold",
		Short: "Scaffolds new resource documentation (Not Implemented)",
		PreRun: func(cmd *cobra.Command, args []string) {
			util.InitLogger(flags.Debug)
			validateProviderDirectoryAccess(fs)
		},
		Run: func(cmd *cobra.Command, args []string) {
			panic("TODO: implement `website-scaffold` functionality")
		},
	}

	rootCmd.AddCommand(validateCmd, fixCmd, scaffoldCmd)

	configureFlags(rootCmd)

	return rootCmd
}

func validateProviderDirectoryAccess(fs afero.Fs) {
	var err error

	if flags.ProviderDirectory == "" {
		flags.ProviderDirectory, err = os.Getwd()
		if err != nil {
			log.WithError(err).Fatal("retrieving current working directory")
		}
	}

	if !util.DirExists(fs, flags.ProviderDirectory) {
		log.WithField("path", flags.ProviderDirectory).Fatal("unable to access provider directory")
	}
}

func printErrors(rd *data.TerraformNodeData) {
	l := len(rd.Errors)
	sep := "---\n\n"

	err := "errors"
	if l == 1 {
		err = "error"
	}

	b := strings.Builder{}
	b.WriteString(util.RedBold(fmt.Sprintf("%s `%s` contains %d %s:\n", rd.Type, rd.Name, l, err)))
	b.WriteString(sep)

	for _, v := range rd.Errors {
		b.WriteString(fmt.Sprintf("-> %s\n", v.Error()))
	}

	b.WriteString("\n")
	b.WriteString(sep)

	fmt.Print(b.String())
}
