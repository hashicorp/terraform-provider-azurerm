// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// CloseEphemeralResourceResponse returns the *tfprotov5.CloseEphemeralResourceResponse
// equivalent of a *fwserver.CloseEphemeralResourceResponse.
func CloseEphemeralResourceResponse(ctx context.Context, fw *fwserver.CloseEphemeralResourceResponse) *tfprotov5.CloseEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.CloseEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
