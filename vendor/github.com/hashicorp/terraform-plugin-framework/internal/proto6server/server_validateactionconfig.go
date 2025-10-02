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

// ValidateActionConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateActionConfig(ctx context.Context, proto6Req *tfprotov6.ValidateActionConfigRequest) (*tfprotov6.ValidateActionConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateActionConfigResponse{}

	action, diags := s.FrameworkServer.Action(ctx, proto6Req.ActionType)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateActionConfigResponse(ctx, fwResp), nil
	}

	actionSchema, diags := s.FrameworkServer.ActionSchema(ctx, proto6Req.ActionType)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateActionConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateActionConfigRequest(ctx, proto6Req, action, actionSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateActionConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateActionConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateActionConfigResponse(ctx, fwResp), nil
}
