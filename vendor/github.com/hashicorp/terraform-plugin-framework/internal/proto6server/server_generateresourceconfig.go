// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
)

// GenerateResourceConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) GenerateResourceConfig(ctx context.Context, proto6Req *tfprotov6.GenerateResourceConfigRequest) (*tfprotov6.GenerateResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.GenerateResourceConfigResponse{}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto6Req.TypeName)
	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.GenerateResourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.GenerateResourceConfigRequest(ctx, proto6Req, resourceSchema)
	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.GenerateResourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.GenerateResourceConfig(ctx, fwReq, fwResp)

	return toproto6.GenerateResourceConfigResponse(ctx, fwResp), nil
}
