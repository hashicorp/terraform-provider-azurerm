package main

import (
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/cmd"
)

func main() {
	c := cmd.Make()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
