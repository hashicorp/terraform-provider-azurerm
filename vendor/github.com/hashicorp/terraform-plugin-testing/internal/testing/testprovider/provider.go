// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

var _ provider.Provider = Provider{}

// Provider is a declarative provider implementation for unit testing in this
// Go module.
type Provider struct {
	Capabilities           *provider.ServerCapabilities
	ConfigureResponse      *provider.ConfigureResponse
	DataSources            map[string]DataSource
	ListResources          map[string]ListResource
	Resources              map[string]Resource
	StateStores            map[string]*StateStore
	SchemaResponse         *provider.SchemaResponse
	StopResponse           *provider.StopResponse
	ValidateConfigResponse *provider.ValidateConfigResponse
}

func (p Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.ConfigureResponse != nil {
		resp.Diagnostics = p.ConfigureResponse.Diagnostics
	}
}

func (p Provider) DataSourcesMap() map[string]datasource.DataSource {
	datasources := make(map[string]datasource.DataSource, len(p.DataSources))

	for typeName, d := range p.DataSources {
		datasources[typeName] = d
	}

	return datasources
}

func (p Provider) ListResourcesMap() map[string]list.ListResource {
	listResources := make(map[string]list.ListResource, len(p.ListResources))

	for typeName, d := range p.ListResources {
		listResources[typeName] = d
	}

	return listResources
}

func (p Provider) ResourcesMap() map[string]resource.Resource {
	resources := make(map[string]resource.Resource, len(p.Resources))

	for typeName, d := range p.Resources {
		resources[typeName] = d
	}

	return resources
}

func (p Provider) StateStoresMap() map[string]statestore.StateStore {
	stateStores := make(map[string]statestore.StateStore, len(p.StateStores))

	for typeName, s := range p.StateStores {
		stateStores[typeName] = s
	}

	return stateStores
}

func (p Provider) Stop(ctx context.Context, req provider.StopRequest, resp *provider.StopResponse) {
	if p.StopResponse != nil {
		resp.Error = p.StopResponse.Error
	}
}

func (p Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	if p.SchemaResponse != nil {
		resp.Diagnostics = p.SchemaResponse.Diagnostics
		resp.Schema = p.SchemaResponse.Schema
	}

	resp.Schema = &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{},
	}
}

func (p Provider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	if p.ValidateConfigResponse != nil {
		resp.Diagnostics = p.ValidateConfigResponse.Diagnostics
	}
}

func (p Provider) ServerCapabilities() *provider.ServerCapabilities {
	if p.Capabilities != nil {
		return p.Capabilities
	}
	return nil
}
