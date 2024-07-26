// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// MoveResourceStateResponse returns the *tfprotov6.MoveResourceStateResponse
// equivalent of a *fwserver.MoveResourceStateResponse.
func MoveResourceStateResponse(ctx context.Context, fw *fwserver.MoveResourceStateResponse) *tfprotov6.MoveResourceStateResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.MoveResourceStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	targetPrivate, diags := fw.TargetPrivate.Bytes(ctx)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.TargetPrivate = targetPrivate

	targetState, diags := State(ctx, fw.TargetState)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.TargetState = targetState

	return proto6
}
