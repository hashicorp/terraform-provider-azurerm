// Copyright (c) HashiCorp, Inc.
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

// ReadResource satisfies the tfprotov5.ProviderServer interface.
func (s *Server) ReadResource(ctx context.Context, proto5Req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ReadResourceResponse{}

	resource, diags := s.FrameworkServer.Resource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ReadResourceResponse(ctx, fwResp), nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ReadResourceResponse(ctx, fwResp), nil
	}

	providerMetaSchema, diags := s.FrameworkServer.ProviderMetaSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ReadResourceResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.ReadResourceRequest(ctx, proto5Req, resource, resourceSchema, providerMetaSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ReadResourceResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ReadResource(ctx, fwReq, fwResp)

	return toproto5.ReadResourceResponse(ctx, fwResp), nil
}
