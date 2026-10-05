// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

type Protov5Provider interface {
	Configure(context.Context, Protov5ConfigureRequest, *Protov5ConfigureResponse)
	DataSourcesMap() map[string]datasource.DataSource
	ResourcesMap() map[string]resource.Resource
	Schema(context.Context, Protov5SchemaRequest, *Protov5SchemaResponse)
	Stop(context.Context, Protov5StopRequest, *Protov5StopResponse)
	ValidateConfig(context.Context, Protov5ValidateConfigRequest, *Protov5ValidateConfigResponse)
}

type Protov5ConfigureRequest struct {
	Config tftypes.Value
}

type Protov5ConfigureResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type Protov5SchemaRequest struct{}

type Protov5SchemaResponse struct {
	Diagnostics []*tfprotov5.Diagnostic
	Schema      *tfprotov5.Schema
}

type Protov5StopRequest struct{}

type Protov5StopResponse struct {
	Error error
}

type Protov5ValidateConfigRequest struct {
	Config tftypes.Value
}

type Protov5ValidateConfigResponse struct {
	Diagnostics []*tfprotov5.Diagnostic
}
