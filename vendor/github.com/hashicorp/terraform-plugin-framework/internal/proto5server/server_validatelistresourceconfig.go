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

// ValidateListResourceConfig satisfies the tfprotov5.ProviderServer interface.
func (s *Server) ValidateListResourceConfig(ctx context.Context, proto5Req *tfprotov5.ValidateListResourceConfigRequest) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateListResourceConfigResponse{}

	listResource, diags := s.FrameworkServer.ListResourceType(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return toproto5.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	listResourceSchema, diags := s.FrameworkServer.ListResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return toproto5.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.ValidateListResourceConfigRequest(ctx, proto5Req, listResource, listResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateListResourceConfig(ctx, fwReq, fwResp)

	return toproto5.ValidateListResourceConfigResponse(ctx, fwResp), nil
}
