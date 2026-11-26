// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// DeleteStateResponse returns the *tfprotov6.DeleteStateResponse
// equivalent of a *fwserver.DeleteStateResponse.
func DeleteStateResponse(ctx context.Context, fw *fwserver.DeleteStateResponse) *tfprotov6.DeleteStateResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.DeleteStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
