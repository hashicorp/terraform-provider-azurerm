// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// ReadResourceResponse returns the *tfprotov5.ReadResourceResponse
// equivalent of a *fwserver.ReadResourceResponse.
func ReadResourceResponse(ctx context.Context, fw *fwserver.ReadResourceResponse) *tfprotov5.ReadResourceResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ReadResourceResponse{
		Deferred:    ResourceDeferred(fw.Deferred),
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	newState, diags := State(ctx, fw.NewState)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.NewState = newState

	newIdentity, diags := ResourceIdentity(ctx, fw.NewIdentity)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.NewIdentity = newIdentity

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Private = newPrivate

	return proto5
}
