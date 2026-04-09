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

// RenewEphemeralResource satisfies the tfprotov5.ProviderServer interface.
func (s *Server) RenewEphemeralResource(ctx context.Context, proto5Req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.RenewEphemeralResourceResponse{}

	ephemeralResource, diags := s.FrameworkServer.EphemeralResource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	ephemeralResourceSchema, diags := s.FrameworkServer.EphemeralResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.RenewEphemeralResourceRequest(ctx, proto5Req, ephemeralResource, ephemeralResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.RenewEphemeralResource(ctx, fwReq, fwResp)

	return toproto5.RenewEphemeralResourceResponse(ctx, fwResp), nil
}
