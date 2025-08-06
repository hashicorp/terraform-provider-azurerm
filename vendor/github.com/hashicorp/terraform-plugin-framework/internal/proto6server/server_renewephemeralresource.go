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

// RenewEphemeralResource satisfies the tfprotov6.ProviderServer interface.
func (s *Server) RenewEphemeralResource(ctx context.Context, proto6Req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.RenewEphemeralResourceResponse{}

	ephemeralResource, diags := s.FrameworkServer.EphemeralResource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	ephemeralResourceSchema, diags := s.FrameworkServer.EphemeralResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.RenewEphemeralResourceRequest(ctx, proto6Req, ephemeralResource, ephemeralResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.RenewEphemeralResourceResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.RenewEphemeralResource(ctx, fwReq, fwResp)

	return toproto6.RenewEphemeralResourceResponse(ctx, fwResp), nil
}
