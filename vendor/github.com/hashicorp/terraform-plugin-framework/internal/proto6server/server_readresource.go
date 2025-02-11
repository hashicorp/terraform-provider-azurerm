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

// ReadResource satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ReadResource(ctx context.Context, proto6Req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ReadResourceResponse{}

	resource, diags := s.FrameworkServer.Resource(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ReadResourceResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ReadResourceResponse(ctx, fwResp), nil
	}

	providerMetaSchema, diags := s.FrameworkServer.ProviderMetaSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ReadResourceResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ReadResourceRequest(ctx, proto6Req, resource, resourceSchema, providerMetaSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ReadResourceResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ReadResource(ctx, fwReq, fwResp)

	return toproto6.ReadResourceResponse(ctx, fwResp), nil
}
