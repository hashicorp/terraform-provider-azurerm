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

// ValidateListResourceConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateListResourceConfig(ctx context.Context, proto6Req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateListResourceConfigResponse{}

	listResource, diags := s.FrameworkServer.ListResourceType(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return toproto6.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	listResourceSchema, diags := s.FrameworkServer.ListResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return toproto6.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateListResourceConfigRequest(ctx, proto6Req, listResource, listResourceSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateListResourceConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateListResourceConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateListResourceConfigResponse(ctx, fwResp), nil
}
