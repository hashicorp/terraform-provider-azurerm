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

// PlanResourceChange satisfies the tfprotov6.ProviderServer interface.
func (s *Server) PlanResourceChange(ctx context.Context, proto6Req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.PlanResourceChangeResponse{}

	resource, diags := s.FrameworkServer.Resource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.PlanResourceChangeResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.PlanResourceChangeResponse(ctx, fwResp), nil
	}

	providerMetaSchema, diags := s.FrameworkServer.ProviderMetaSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.PlanResourceChangeResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.PlanResourceChangeRequest(ctx, proto6Req, resource, resourceSchema, providerMetaSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.PlanResourceChangeResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.PlanResourceChange(ctx, fwReq, fwResp)

	return toproto6.PlanResourceChangeResponse(ctx, fwResp), nil
}
