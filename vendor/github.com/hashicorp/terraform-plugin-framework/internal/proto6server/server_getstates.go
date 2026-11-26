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

func (s *Server) GetStates(ctx context.Context, proto6Req *tfprotov6.GetStatesRequest) (*tfprotov6.GetStatesResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.GetStatesResponse{}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.GetStatesResponse(ctx, fwResp), nil
	}

	fwReq := &fwserver.GetStatesRequest{
		StateStore: stateStore,
	}

	s.FrameworkServer.GetStates(ctx, fwReq, fwResp)

	return toproto6.GetStatesResponse(ctx, fwResp), nil
}
