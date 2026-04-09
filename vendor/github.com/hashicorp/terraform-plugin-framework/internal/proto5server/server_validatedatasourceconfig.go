// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateDataSourceConfig satisfies the tfprotov5.ProviderServer interface.
func (s *Server) ValidateDataSourceConfig(ctx context.Context, proto5Req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateDataSourceConfigResponse{}

	dataSource, diags := s.FrameworkServer.DataSource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	dataSourceSchema, diags := s.FrameworkServer.DataSourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.ValidateDataSourceConfigRequest(ctx, proto5Req, dataSource, dataSourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ValidateDataSourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateDataSourceConfig(ctx, fwReq, fwResp)

	return toproto5.ValidateDataSourceConfigResponse(ctx, fwResp), nil
}
