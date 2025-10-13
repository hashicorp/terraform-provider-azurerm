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

// UpgradeResourceIdentity satisfies the tfprotov5.ProviderServer interface.
func (s *Server) UpgradeResourceIdentity(ctx context.Context, proto5Req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.UpgradeResourceIdentityResponse{}

	if proto5Req == nil {
		return toproto5.UpgradeResourceIdentityResponse(ctx, fwResp), nil
	}

	resource, diags := s.FrameworkServer.Resource(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceIdentityResponse(ctx, fwResp), nil
	}

	identitySchema, diags := s.FrameworkServer.ResourceIdentitySchema(ctx, proto5Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceIdentityResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.UpgradeResourceIdentityRequest(ctx, proto5Req, resource, identitySchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.UpgradeResourceIdentityResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.UpgradeResourceIdentity(ctx, fwReq, fwResp)

	return toproto5.UpgradeResourceIdentityResponse(ctx, fwResp), nil
}
