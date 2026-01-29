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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
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
func ProtoV5ProviderFactoriesInitWithHTTPClient(ctx context.Context, httpClient *http.Client, providerNames ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, _, err := ProtoV5ProviderServerFactoryWithHTTPClient(ctx, httpClient)
			if err != nil {
				return nil, err
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

// ProtoV5ProviderServerFactoryWithHTTPClient creates a provider server factory with a custom HTTP client.
func ProtoV5ProviderServerFactoryWithHTTPClient(ctx context.Context, httpClient *http.Client) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	v2Provider := provider.AzureProviderWithHTTPClient(httpClient)

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
