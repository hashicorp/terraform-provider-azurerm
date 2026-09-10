// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var _ provider.Protov5Provider = Protov5Provider{}

// Protov5Provider is a declarative provider implementation for unit testing in this
// Go module. The provider is unimplemented except for the Schema method.
type Protov5Provider struct {
	ConfigureResponse      *provider.Protov5ConfigureResponse
	DataSources            map[string]DataSource
	Resources              map[string]Resource
	SchemaResponse         *provider.Protov5SchemaResponse
	StopResponse           *provider.Protov5StopResponse
	ValidateConfigResponse *provider.Protov5ValidateConfigResponse
}

func (p Protov5Provider) Configure(ctx context.Context, req provider.Protov5ConfigureRequest, resp *provider.Protov5ConfigureResponse) {

}

func (p Protov5Provider) DataSourcesMap() map[string]datasource.DataSource {
	return nil
}

func (p Protov5Provider) ResourcesMap() map[string]resource.Resource {
	return nil
}

func (p Protov5Provider) Stop(ctx context.Context, req provider.Protov5StopRequest, resp *provider.Protov5StopResponse) {
}

func (p Protov5Provider) Schema(ctx context.Context, req provider.Protov5SchemaRequest, resp *provider.Protov5SchemaResponse) {
	resp.Schema = &tfprotov5.Schema{
		Block: &tfprotov5.SchemaBlock{},
	}
}

func (p Protov5Provider) ValidateConfig(ctx context.Context, req provider.Protov5ValidateConfigRequest, resp *provider.Protov5ValidateConfigResponse) {
}
