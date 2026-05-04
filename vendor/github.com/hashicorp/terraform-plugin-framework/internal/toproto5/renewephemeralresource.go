// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// RenewEphemeralResourceResponse returns the *tfprotov5.RenewEphemeralResourceResponse
// equivalent of a *fwserver.RenewEphemeralResourceResponse.
func RenewEphemeralResourceResponse(ctx context.Context, fw *fwserver.RenewEphemeralResourceResponse) *tfprotov5.RenewEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.RenewEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		RenewAt:     fw.RenewAt,
	}

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Private = newPrivate

	return proto5
}
