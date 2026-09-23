// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// LockStateResponse returns the *tfprotov6.LockStateResponse
// equivalent of a *fwserver.LockStateResponse.
func LockStateResponse(ctx context.Context, fw *fwserver.LockStateResponse) *tfprotov6.LockStateResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.LockStateResponse{
		LockID:      fw.LockID,
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
