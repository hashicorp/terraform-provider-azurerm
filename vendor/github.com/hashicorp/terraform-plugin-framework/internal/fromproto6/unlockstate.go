// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// UnlockStateRequest returns the *fwserver.UnlockStateRequest
// equivalent of a *tfprotov6.UnlockStateRequest.
func UnlockStateRequest(ctx context.Context, proto6 *tfprotov6.UnlockStateRequest, reqStateStore statestore.StateStore) *fwserver.UnlockStateRequest {
	if proto6 == nil {
		return nil
	}

	return &fwserver.UnlockStateRequest{
		StateID:    proto6.StateID,
		LockID:     proto6.LockID,
		StateStore: reqStateStore,
	}
}
