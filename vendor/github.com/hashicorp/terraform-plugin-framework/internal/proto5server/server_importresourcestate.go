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

// ImportResourceState satisfies the tfprotov5.ProviderServer interface.
func (s *Server) ImportResourceState(ctx context.Context, proto5Req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ImportResourceStateResponse{}

	resource, diags := s.FrameworkServer.Resource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ImportResourceStateResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ImportResourceStateResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.ImportResourceStateRequest(ctx, proto5Req, resource, resourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ImportResourceStateResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ImportResourceState(ctx, fwReq, fwResp)

	return toproto5.ImportResourceStateResponse(ctx, fwResp), nil
}
