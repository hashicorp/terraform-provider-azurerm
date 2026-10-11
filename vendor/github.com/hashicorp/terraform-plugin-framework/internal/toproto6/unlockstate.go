// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// UnlockStateResponse returns the *tfprotov6.UnlockStateResponse
// equivalent of a *fwserver.UnlockStateResponse.
func UnlockStateResponse(ctx context.Context, fw *fwserver.UnlockStateResponse) *tfprotov6.UnlockStateResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.UnlockStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
