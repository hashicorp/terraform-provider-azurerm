// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// WriteStateBytesResponse returns the *tfprotov6.WriteStateBytesResponse
// equivalent of a *fwserver.WriteStateBytesResponse.
func WriteStateBytesResponse(ctx context.Context, fw *fwserver.WriteStateBytesResponse) *tfprotov6.WriteStateBytesResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.WriteStateBytesResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
