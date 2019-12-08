package main

import (
	"log"

	"github.com/Azure/go-autorest/tracing"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/tracer"
)

func main() {
	tracing.Register(tracer.MyTracer)

	// enable tracer
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider})
}
