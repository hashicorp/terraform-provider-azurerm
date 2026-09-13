// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

type DeleteStateRequest struct {
	StateID    string
	StateStore statestore.StateStore
}

type DeleteStateResponse struct {
	Diagnostics diag.Diagnostics
}

// DeleteState implements the framework server DeleteState RPC.
func (s *Server) DeleteState(ctx context.Context, req *DeleteStateRequest, resp *DeleteStateResponse) {
	if req == nil {
		return
	}

	if stateStoreWithConfigure, ok := req.StateStore.(statestore.StateStoreWithConfigure); ok {
		logging.FrameworkTrace(ctx, "StateStore implements StateStoreWithConfigure")

		configureReq := statestore.ConfigureRequest{
			StateStoreData: s.StateStoreConfigureData.StateStoreConfigureData,
		}
		configureResp := statestore.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined StateStore Configure")
		stateStoreWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined StateStore Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	deleteStateReq := statestore.DeleteStateRequest{
		StateID: req.StateID,
	}
	deleteStateResp := statestore.DeleteStateResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore DeleteState")
	req.StateStore.DeleteState(ctx, deleteStateReq, &deleteStateResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore DeleteState")

	resp.Diagnostics.Append(deleteStateResp.Diagnostics...)
}
