package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider})
}
