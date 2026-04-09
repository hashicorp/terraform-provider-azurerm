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

// ValidateEphemeralResourceConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateEphemeralResourceConfig(ctx context.Context, proto6Req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateEphemeralResourceConfigResponse{}

	ephemeralResource, diags := s.FrameworkServer.EphemeralResource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateEphemeralResourceConfigResponse(ctx, fwResp), nil
	}

	ephemeralResourceSchema, diags := s.FrameworkServer.EphemeralResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateEphemeralResourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateEphemeralResourceConfigRequest(ctx, proto6Req, ephemeralResource, ephemeralResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateEphemeralResourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateEphemeralResourceConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateEphemeralResourceConfigResponse(ctx, fwResp), nil
}
