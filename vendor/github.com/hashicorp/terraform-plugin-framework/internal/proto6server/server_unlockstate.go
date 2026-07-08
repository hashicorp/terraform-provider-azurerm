// Copyright IBM Corp. 2021, 2026
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

func (s *Server) UnlockState(ctx context.Context, req *tfprotov6.UnlockStateRequest) (*tfprotov6.UnlockStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.UnlockStateResponse{}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.UnlockStateResponse(ctx, fwResp), nil
	}

	fwReq := fromproto6.UnlockStateRequest(ctx, req, stateStore)

	s.FrameworkServer.UnlockState(ctx, fwReq, fwResp)

	return toproto6.UnlockStateResponse(ctx, fwResp), nil
}
