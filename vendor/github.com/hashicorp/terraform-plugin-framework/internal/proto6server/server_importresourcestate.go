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

// ImportResourceState satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ImportResourceState(ctx context.Context, proto6Req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ImportResourceStateResponse{}

	resource, diags := s.FrameworkServer.Resource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ImportResourceStateResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ImportResourceStateResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ImportResourceStateRequest(ctx, proto6Req, resource, resourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ImportResourceStateResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ImportResourceState(ctx, fwReq, fwResp)

	return toproto6.ImportResourceStateResponse(ctx, fwResp), nil
}
