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

func (s *Server) ConfigureStateStore(ctx context.Context, req *tfprotov6.ConfigureStateStoreRequest) (*tfprotov6.ConfigureStateStoreResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ConfigureStateStoreResponse{}

	stateStore, diags := s.FrameworkServer.StateStore(ctx, req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ConfigureStateStoreResponse(ctx, fwResp), nil
	}

	stateStoreSchema, diags := s.FrameworkServer.StateStoreSchema(ctx, req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ConfigureStateStoreResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ConfigureStateStoreRequest(ctx, req, stateStore, stateStoreSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ConfigureStateStoreResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ConfigureStateStore(ctx, fwReq, fwResp)

	return toproto6.ConfigureStateStoreResponse(ctx, fwResp), nil
}
