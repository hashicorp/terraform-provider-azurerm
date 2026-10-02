// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

type LockStateRequest struct {
	StateID    string
	Operation  string
	StateStore statestore.StateStore
}

type LockStateResponse struct {
	LockID      string
	Diagnostics diag.Diagnostics
}

func (s *Server) LockState(ctx context.Context, req *LockStateRequest, resp *LockStateResponse) {
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

	lockReq := statestore.LockRequest{
		StateID:   req.StateID,
		Operation: req.Operation,
	}

	lockResp := statestore.LockResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore Lock")
	req.StateStore.Lock(ctx, lockReq, &lockResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore Lock")

	resp.Diagnostics.Append(lockResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.LockID = lockResp.LockID
}
