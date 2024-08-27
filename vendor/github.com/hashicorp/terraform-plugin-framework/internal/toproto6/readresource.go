// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// ReadResourceResponse returns the *tfprotov6.ReadResourceResponse
// equivalent of a *fwserver.ReadResourceResponse.
func ReadResourceResponse(ctx context.Context, fw *fwserver.ReadResourceResponse) *tfprotov6.ReadResourceResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ReadResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	newState, diags := State(ctx, fw.NewState)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.NewState = newState

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.Private = newPrivate

	return proto6
}
