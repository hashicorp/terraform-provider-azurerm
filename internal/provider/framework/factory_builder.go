// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

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
