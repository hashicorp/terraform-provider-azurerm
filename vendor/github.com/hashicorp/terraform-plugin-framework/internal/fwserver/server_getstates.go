// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

type GetStatesRequest struct {
	StateStore statestore.StateStore
}

type GetStatesResponse struct {
	StateIDs    []string
	Diagnostics diag.Diagnostics
}

// GetStates implements the framework server GetStates RPC.
func (s *Server) GetStates(ctx context.Context, req *GetStatesRequest, resp *GetStatesResponse) {
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

	getStatesReq := statestore.GetStatesRequest{}
	getStatesResp := statestore.GetStatesResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore GetStates")
	req.StateStore.GetStates(ctx, getStatesReq, &getStatesResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore GetStates")

	resp.Diagnostics.Append(getStatesResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.StateIDs = getStatesResp.StateIDs
}
