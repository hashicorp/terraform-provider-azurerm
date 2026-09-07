// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
)

// GenerateResourceConfig satisfies the tfprotov5.ProviderServer interface.
func (s *Server) GenerateResourceConfig(ctx context.Context, proto5Req *tfprotov5.GenerateResourceConfigRequest) (*tfprotov5.GenerateResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.GenerateResourceConfigResponse{}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto5Req.TypeName)
	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.GenerateResourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.GenerateResourceConfigRequest(ctx, proto5Req, resourceSchema)
	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.GenerateResourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.GenerateResourceConfig(ctx, fwReq, fwResp)

	return toproto5.GenerateResourceConfigResponse(ctx, fwResp), nil
}
