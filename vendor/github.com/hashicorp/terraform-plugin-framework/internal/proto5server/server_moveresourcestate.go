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

// MoveResourceState satisfies the tfprotov5.ProviderServer interface.
func (s *Server) MoveResourceState(ctx context.Context, proto5Req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.MoveResourceStateResponse{}

	if proto5Req == nil {
		return toproto5.MoveResourceStateResponse(ctx, fwResp), nil
	}

	resource, diags := s.FrameworkServer.Resource(ctx, proto5Req.TargetTypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.MoveResourceStateResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto5Req.TargetTypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.MoveResourceStateResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.MoveResourceStateRequest(ctx, proto5Req, resource, resourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.MoveResourceStateResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.MoveResourceState(ctx, fwReq, fwResp)

	return toproto5.MoveResourceStateResponse(ctx, fwResp), nil
}
