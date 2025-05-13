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
	vRulesFlag    string
	vServiceFlag  string
	vResourceFlag string

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates resource documentation",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			util.InitLogger(debugFlag)
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
			if vRulesFlag != "" {
				rules = strings.Split(vRulesFlag, ",")
			}

			resources := data.GetData(fs, providerDirFlag, vServiceFlag, vResourceFlag)

			v := validator.Validator{}
			for _, r := range resources {
				v.Run(rules, r, false)
			}

			// separate from loop above to prevent debug messages from mixing in
			errCount, resourceWithErrCount := 0, 0
			for _, r := range resources {
				if l := len(r.Errors); l > 0 {
					resourceWithErrCount++
					errCount += l
					r.PrintErrors()
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
			} else {
				fmt.Print(util.GreenBold("Found no errors\n"))
			}
		},
	}
)

func init() {
	validateCmd.Flags().StringVar(&vRulesFlag, "rules", "", "comma separated list of rule names, if not specified, all rules will be run")
	validateCmd.Flags().StringVarP(&vServiceFlag, "service", "s", "", "service to validate")
	validateCmd.Flags().StringVarP(&vResourceFlag, "resource", "r", "", "resource to validate")

	rootCmd.AddCommand(validateCmd)
}
