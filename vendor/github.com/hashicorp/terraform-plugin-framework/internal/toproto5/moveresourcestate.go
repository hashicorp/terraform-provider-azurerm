// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// MoveResourceStateResponse returns the *tfprotov5.MoveResourceStateResponse
// equivalent of a *fwserver.MoveResourceStateResponse.
func MoveResourceStateResponse(ctx context.Context, fw *fwserver.MoveResourceStateResponse) *tfprotov5.MoveResourceStateResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.MoveResourceStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	targetPrivate, diags := fw.TargetPrivate.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.TargetPrivate = targetPrivate

	targetState, diags := State(ctx, fw.TargetState)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.TargetState = targetState

	return proto5
}
