// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// CloseEphemeralResourceResponse returns the *tfprotov6.CloseEphemeralResourceResponse
// equivalent of a *fwserver.CloseEphemeralResourceResponse.
func CloseEphemeralResourceResponse(ctx context.Context, fw *fwserver.CloseEphemeralResourceResponse) *tfprotov6.CloseEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.CloseEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
