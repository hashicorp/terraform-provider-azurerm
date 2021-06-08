package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		// This is needed so Terraform doesn't default to expecting protocol 4.
		// TODO: remove below line once the provider migrates to plugin SDK v2.
		os.Setenv("PLUGIN_PROTOCOL_VERSIONS", "5")
		err := plugin.Debug(context.Background(), "registry.terraform.io/hashicorp/azurerm",
			&plugin.ServeOpts{
				ProviderFunc: azurerm.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: azurerm.Provider,
		})
	}
}
