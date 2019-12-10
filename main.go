package main

import (
	"context"
	"log"

	opencensusTrace "go.opencensus.io/trace"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/tracer"
)

func main() {
	// enable tracer
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if tracer.TracingEnabled() {
		tracer.Init()
		// create the first root span, this span has the same lifetime as the plugin server
		_, tracer.RootSpan = opencensusTrace.StartSpan(context.Background(), "root")
		defer tracer.RootSpan.End()
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider})
}
