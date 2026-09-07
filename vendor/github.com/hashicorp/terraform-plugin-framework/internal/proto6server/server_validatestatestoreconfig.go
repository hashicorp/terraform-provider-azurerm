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

// ValidateStateStoreConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateStateStoreConfig(ctx context.Context, proto6Req *tfprotov6.ValidateStateStoreConfigRequest) (*tfprotov6.ValidateStateStoreConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateStateStoreConfigResponse{}

	statestore, diags := s.FrameworkServer.StateStore(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateStateStoreConfigResponse(ctx, fwResp), nil
	}

	statestoreSchema, diags := s.FrameworkServer.StateStoreSchema(ctx, proto6Req.TypeName)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateStateStoreConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateStateStoreConfigRequest(ctx, proto6Req, statestore, statestoreSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateStateStoreConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateStateStoreConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateStateStoreConfigResponse(ctx, fwResp), nil
}
