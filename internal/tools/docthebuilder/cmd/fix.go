package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	fRulesFlag    string
	fServiceFlag  string
	fResourceFlag string

	fixCmd = &cobra.Command{
		Use:   "fix",
		Short: "Attempts to fix resource documentation",
		Run: func(cmd *cobra.Command, args []string) {
			util.InitLogger(debugFlag)

			var err error

			fs := afero.NewOsFs()

			if providerDirFlag == "" {
				providerDirFlag, err = os.Getwd()
				if err != nil {
					log.WithError(err).Fatal("retrieving current working directory")
				}
			}

			if !util.DirExists(fs, providerDirFlag) {
				log.WithField("path", providerDirFlag).Fatal("unable to access provider directory")
			}

			rules := util.MapKeys2Slice(rule.Registration)
			if fRulesFlag != "" {
				rules = strings.Split(fRulesFlag, ",")
			}

			resources := data.GetData(fs, providerDirFlag, fServiceFlag, fResourceFlag)

			v := validator.Validator{}
			errCount, resourceWithErrCount := 0, 0
			for _, r := range resources {
				v.Run(rules, r, true)

				if r.Document.HasChange {
					resourceWithErrCount++
					errCount = errCount + len(r.Errors)

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
					errStr = errStr + "s"
				}

				if resourceWithErrCount > 1 {
					resourceStr = resourceStr + "s"
				}
				fmt.Printf(util.Red("Found %d %s in %d %s and applied fixes where possible, please review the changes\n"), errCount, errStr, resourceWithErrCount, resourceStr)
			} else {
				fmt.Print(util.GreenBold("Found no errors\n"))
			}
		},
	}
)

func init() {
	fixCmd.Flags().StringVar(&fRulesFlag, "rules", "", "comma separated list of rule names, if not specified, all rules will be run")
	fixCmd.Flags().StringVarP(&fServiceFlag, "service", "s", "", "service to validate and fix")
	fixCmd.Flags().StringVarP(&fResourceFlag, "resource", "r", "", "resource to validate and fix")

	rootCmd.AddCommand(fixCmd)
}
