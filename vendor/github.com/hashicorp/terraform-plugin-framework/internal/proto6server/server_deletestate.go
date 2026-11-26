// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func (s *Server) DeleteState(ctx context.Context, proto6Req *tfprotov6.DeleteStateRequest) (*tfprotov6.DeleteStateResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.DeleteStateResponse{}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.DeleteStateResponse(ctx, fwResp), nil
	}

	fwReq := &fwserver.DeleteStateRequest{
		StateID:    proto6Req.StateID,
		StateStore: stateStore,
	}

	s.FrameworkServer.DeleteState(ctx, fwReq, fwResp)

	return toproto6.DeleteStateResponse(ctx, fwResp), nil
}
