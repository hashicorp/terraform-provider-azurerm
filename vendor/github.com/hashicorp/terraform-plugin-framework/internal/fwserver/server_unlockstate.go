// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

type UnlockStateRequest struct {
	StateID    string
	LockID     string
	StateStore statestore.StateStore
}

type UnlockStateResponse struct {
	Diagnostics diag.Diagnostics
}

func (s *Server) UnlockState(ctx context.Context, req *UnlockStateRequest, resp *UnlockStateResponse) {
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

	unlockReq := statestore.UnlockRequest{
		StateID: req.StateID,
		LockID:  req.LockID,
	}

	unlockResp := statestore.UnlockResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore Unlock")
	req.StateStore.Unlock(ctx, unlockReq, &unlockResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore Unlock")

	resp.Diagnostics.Append(unlockResp.Diagnostics...)
}
