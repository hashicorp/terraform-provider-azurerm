// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetStatesResponse returns the *tfprotov6.GetStatesResponse
// equivalent of a *fwserver.GetStatesResponse.
func GetStatesResponse(ctx context.Context, fw *fwserver.GetStatesResponse) *tfprotov6.GetStatesResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.GetStatesResponse{
		StateIDs:    fw.StateIDs,
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
