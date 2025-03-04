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

// MoveResourceState satisfies the tfprotov6.ProviderServer interface.
func (s *Server) MoveResourceState(ctx context.Context, proto6Req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.MoveResourceStateResponse{}

	if proto6Req == nil {
		return toproto6.MoveResourceStateResponse(ctx, fwResp), nil
	}

	resource, diags := s.FrameworkServer.Resource(ctx, proto6Req.TargetTypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.MoveResourceStateResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto6Req.TargetTypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.MoveResourceStateResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.MoveResourceStateRequest(ctx, proto6Req, resource, resourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.MoveResourceStateResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.MoveResourceState(ctx, fwReq, fwResp)

	return toproto6.MoveResourceStateResponse(ctx, fwResp), nil
}
