package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	debugFlag       bool
	providerDirFlag string

	rootCmd = &cobra.Command{
		Use:   "docthebuilder [validate|fix|scaffold]",
		Short: "TODO",
		Long:  `TODO`,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "output debug logs, defaults to false")
	rootCmd.PersistentFlags().StringVarP(&providerDirFlag, "provider-directory", "p", "", "provider directory path, can be omitted if provider directory is the current working directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
