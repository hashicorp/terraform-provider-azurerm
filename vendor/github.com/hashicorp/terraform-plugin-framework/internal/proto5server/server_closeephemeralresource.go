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

// CloseEphemeralResource satisfies the tfprotov5.ProviderServer interface.
func (s *Server) CloseEphemeralResource(ctx context.Context, proto5Req *tfprotov5.CloseEphemeralResourceRequest) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.CloseEphemeralResourceResponse{}

	ephemeralResource, diags := s.FrameworkServer.EphemeralResource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.CloseEphemeralResourceResponse(ctx, fwResp), nil
	}

	ephemeralResourceSchema, diags := s.FrameworkServer.EphemeralResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.CloseEphemeralResourceResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.CloseEphemeralResourceRequest(ctx, proto5Req, ephemeralResource, ephemeralResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.CloseEphemeralResourceResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.CloseEphemeralResource(ctx, fwReq, fwResp)

	return toproto5.CloseEphemeralResourceResponse(ctx, fwResp), nil
}
