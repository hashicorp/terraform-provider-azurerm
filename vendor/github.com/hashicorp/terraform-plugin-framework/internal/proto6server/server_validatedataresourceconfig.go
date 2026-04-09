// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateDataResourceConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateDataResourceConfig(ctx context.Context, proto6Req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateDataSourceConfigResponse{}

	dataSource, diags := s.FrameworkServer.DataSource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	dataSourceSchema, diags := s.FrameworkServer.DataSourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateDataSourceConfigRequest(ctx, proto6Req, dataSource, dataSourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateDataSourceConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateDataSourceConfigResponse(ctx, fwResp), nil
}
