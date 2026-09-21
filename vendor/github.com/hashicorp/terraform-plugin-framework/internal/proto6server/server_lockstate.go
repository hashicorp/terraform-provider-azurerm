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

func (s *Server) LockState(ctx context.Context, req *tfprotov6.LockStateRequest) (*tfprotov6.LockStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.LockStateResponse{}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.LockStateResponse(ctx, fwResp), nil
	}

	fwReq := fromproto6.LockStateRequest(ctx, req, stateStore)

	s.FrameworkServer.LockState(ctx, fwReq, fwResp)

	return toproto6.LockStateResponse(ctx, fwResp), nil
}
