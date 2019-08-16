package main

import (
	"log"

	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
)

func main() {

	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider})
}
