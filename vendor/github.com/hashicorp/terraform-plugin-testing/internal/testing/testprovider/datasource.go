// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
)

var _ datasource.DataSource = DataSource{}

type DataSource struct {
	ReadResponse           *datasource.ReadResponse
	SchemaResponse         *datasource.SchemaResponse
	ValidateConfigResponse *datasource.ValidateConfigResponse
}

func (d DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.ReadResponse != nil {
		resp.Diagnostics = d.ReadResponse.Diagnostics
		resp.State = d.ReadResponse.State
	}
}

func (d DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	if d.SchemaResponse != nil {
		resp.Diagnostics = d.SchemaResponse.Diagnostics
		resp.Schema = d.SchemaResponse.Schema
	}
}

func (d DataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if d.ValidateConfigResponse != nil {
		resp.Diagnostics = d.ValidateConfigResponse.Diagnostics
	}
}
