// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	providerServer, _, err := framework.ProtoV5ProviderServerFactory(ctx)
	if err != nil {
		log.Fatalf("creating AzureRM Provider Server: %+v", err)
	}

	var serveOpts []tf5server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve("registry.terraform.io/hashicorp/azurerm", providerServer, serveOpts...)
	if err != nil {
		log.Fatal(err)
	}
}
