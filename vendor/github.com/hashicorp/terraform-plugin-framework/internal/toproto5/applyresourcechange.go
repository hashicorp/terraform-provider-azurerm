// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// ApplyResourceChangeResponse returns the *tfprotov5.ApplyResourceChangeResponse
// equivalent of a *fwserver.ApplyResourceChangeResponse.
func ApplyResourceChangeResponse(ctx context.Context, fw *fwserver.ApplyResourceChangeResponse) *tfprotov5.ApplyResourceChangeResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ApplyResourceChangeResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	newState, diags := State(ctx, fw.NewState)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.NewState = newState

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Private = newPrivate

	return proto5
}
