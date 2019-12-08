package main

import (
	"context"
	"log"

	"github.com/Azure/go-autorest/tracing"

	"go.opencensus.io/trace"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/tracer"
)

func main() {
	tracing.Register(tracer.Tracer)
	_, tracer.RootSpan = trace.StartSpan(context.Background(), "root")
	defer tracer.RootSpan.End()

	// enable tracer
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider})
}
