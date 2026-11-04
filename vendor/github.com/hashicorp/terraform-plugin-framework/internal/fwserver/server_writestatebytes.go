// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

type WriteStateBytesRequest struct {
	StateID    string
	StateBytes []byte
	StateStore statestore.StateStore
}

type WriteStateBytesResponse struct {
	Diagnostics diag.Diagnostics
}

// WriteStateBytes implements the framework server WriteStateBytes RPC.
func (s *Server) WriteStateBytes(ctx context.Context, req *WriteStateBytesRequest, resp *WriteStateBytesResponse) {
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

	writeReq := statestore.WriteRequest{
		StateID:    req.StateID,
		StateBytes: req.StateBytes,
	}

	writeResp := statestore.WriteResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore Write")
	req.StateStore.Write(ctx, writeReq, &writeResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore Write")

	resp.Diagnostics.Append(writeResp.Diagnostics...)
}
