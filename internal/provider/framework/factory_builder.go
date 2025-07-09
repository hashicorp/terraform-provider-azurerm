// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"

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

// Refactor: move to terraform-plugin-sdk
type SDKContext string

// Refactor: move to terraform-plugin-sdk
var SDKResourceKey SDKContext = "sdk_resource"

// Refactor: move to terraform-plugin-sdk
// NewContextWithSDKResource returns a new Context that carries value r
func NewContextWithSDKResource(ctx context.Context, r *schema.Resource) context.Context {
	return context.WithValue(ctx, SDKResourceKey, r)
}

// Refactor: move to terraform-plugin-sdk
// FromContext returns the SDK Resource value stored in ctx, if any.
func SDKResourceFromContext(ctx context.Context) (*schema.Resource, bool) {
	r, ok := ctx.Value(SDKResourceKey).(*schema.Resource)
	return r, ok
}

func ProtoV5ProviderServerFactory(ctx context.Context) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	v2Provider := provider.AzureProvider()

	providers := []func() tfprotov5.ProviderServer{
		v2Provider.GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider(v2Provider)),
	}

	listInterceptor := func(ctx context.Context, req *tfprotov5.ListResourceRequest) context.Context {
		typeName := req.TypeName
		resource, ok := v2Provider.ResourcesMap[typeName]
		if !ok {
			return ctx
		}

		return NewContextWithSDKResource(ctx, resource)
	}
	interceptor := tf5muxserver.Interceptor{BeforeListResource: listInterceptor}

	muxServer, err := tf5muxserver.NewMuxServerWithOptions(ctx, Servers(providers...), Interceptors(interceptor))
	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, v2Provider, nil
}

func V5ProviderWithoutPluginSDK() func() tfprotov5.ProviderServer {
	return providerserver.NewProtocol5(NewFrameworkV5Provider())
}
