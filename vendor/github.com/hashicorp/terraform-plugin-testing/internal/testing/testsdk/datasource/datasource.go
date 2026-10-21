// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type DataSource interface {
	Read(context.Context, ReadRequest, *ReadResponse)
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}

type ReadRequest struct {
	Config tftypes.Value
}

type ReadResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	State       tftypes.Value
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}

type ValidateConfigRequest struct {
	Config tftypes.Value
}

type ValidateConfigResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}
