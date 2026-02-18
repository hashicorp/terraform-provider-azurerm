// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

func ProtoV6ProviderFactoriesInit(_ context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {
		// This is all we need from protoV6 for now to properly test ephemeral resources
		if name == "echo" {
			factories[name] = echoprovider.NewProviderServer()
		}
	}

	return factories
}

func ProtoV5ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, _, err := ProtoV5ProviderServerFactory(ctx)
			if err != nil {
				return nil, err
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

// ProtoV5ProviderFactoriesInitWithHTTPClient creates provider factories with a custom HTTP client.
// The HTTP client is injected into provider.Meta() before ConfigureContextFunc runs,
// allowing it to be used by the Azure SDK during client initialization.
func ProtoV5ProviderFactoriesInitWithHTTPClient(ctx context.Context, httpClient *http.Client, providerNames ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, v2Provider, err := ProtoV5ProviderServerFactory(ctx)
			if err != nil {
				return nil, err
			}

			// Wrap the original ConfigureContextFunc to inject HTTPClient into Meta
			configureContextFunc := v2Provider.ConfigureContextFunc
			v2Provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				// Set HTTPClient in Meta before ConfigureContextFunc runs
				var meta *clients.Client
				if v, ok := v2Provider.Meta().(*clients.Client); ok && v != nil {
					meta = v
				} else {
					meta = new(clients.Client)
				}
				meta.HTTPClient = httpClient
				v2Provider.SetMeta(meta)

				// Call the original ConfigureContextFunc
				return configureContextFunc(ctx, d)
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

func ProtoV5ProviderServerFactory(ctx context.Context) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	v2Provider := provider.AzureProvider()

	providers := []func() tfprotov5.ProviderServer{
		v2Provider.GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider(v2Provider)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, v2Provider, nil
}

func V5ProviderWithoutPluginSDK() func() tfprotov5.ProviderServer {
	return providerserver.NewProtocol5(NewFrameworkV5Provider())
}
