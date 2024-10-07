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

// UpgradeResourceState satisfies the tfprotov5.ProviderServer interface.
func (s *Server) UpgradeResourceState(ctx context.Context, proto5Req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.UpgradeResourceStateResponse{}

	if proto5Req == nil {
		return toproto5.UpgradeResourceStateResponse(ctx, fwResp), nil
	}

	resource, diags := s.FrameworkServer.Resource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceStateResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceStateResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.UpgradeResourceStateRequest(ctx, proto5Req, resource, resourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceStateResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.UpgradeResourceState(ctx, fwReq, fwResp)

	return toproto5.UpgradeResourceStateResponse(ctx, fwResp), nil
}
