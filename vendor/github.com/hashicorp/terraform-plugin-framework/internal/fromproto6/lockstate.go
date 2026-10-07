// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// LockStateRequest returns the *fwserver.LockStateRequest
// equivalent of a *tfprotov6.LockStateRequest.
func LockStateRequest(ctx context.Context, proto6 *tfprotov6.LockStateRequest, reqStateStore statestore.StateStore) *fwserver.LockStateRequest {
	if proto6 == nil {
		return nil
	}

	return &fwserver.LockStateRequest{
		StateID:    proto6.StateID,
		Operation:  proto6.Operation,
		StateStore: reqStateStore,
	}
}
