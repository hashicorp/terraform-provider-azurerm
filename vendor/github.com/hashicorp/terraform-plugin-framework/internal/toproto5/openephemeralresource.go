// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// OpenEphemeralResourceResponse returns the *tfprotov5.OpenEphemeralResourceResponse
// equivalent of a *fwserver.OpenEphemeralResourceResponse.
func OpenEphemeralResourceResponse(ctx context.Context, fw *fwserver.OpenEphemeralResourceResponse) *tfprotov5.OpenEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.OpenEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		RenewAt:     fw.RenewAt,
		Deferred:    EphemeralResourceDeferred(fw.Deferred),
	}

	result, diags := EphemeralResultData(ctx, fw.Result)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Result = result

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Private = newPrivate

	return proto5
}
