// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// RenewEphemeralResourceResponse returns the *tfprotov6.RenewEphemeralResourceResponse
// equivalent of a *fwserver.RenewEphemeralResourceResponse.
func RenewEphemeralResourceResponse(ctx context.Context, fw *fwserver.RenewEphemeralResourceResponse) *tfprotov6.RenewEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.RenewEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		RenewAt:     fw.RenewAt,
	}

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.Private = newPrivate

	return proto6
}
