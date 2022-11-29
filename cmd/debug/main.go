package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	opts := &plugin.ServeOpts{
		Debug:        true,
		ProviderFunc: provider.AzureProvider,
		ProviderAddr: "snyk/azurerm",
	}
	plugin.Serve(opts)
}
