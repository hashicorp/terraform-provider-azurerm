package cmd

import "github.com/spf13/cobra"

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffolds new resource documentation (Not Implemented)",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		panic("TODO: implement `website-scaffold` functionality")
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}
