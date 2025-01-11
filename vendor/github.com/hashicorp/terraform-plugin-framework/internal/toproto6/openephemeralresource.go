// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// OpenEphemeralResourceResponse returns the *tfprotov6.OpenEphemeralResourceResponse
// equivalent of a *fwserver.OpenEphemeralResourceResponse.
func OpenEphemeralResourceResponse(ctx context.Context, fw *fwserver.OpenEphemeralResourceResponse) *tfprotov6.OpenEphemeralResourceResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.OpenEphemeralResourceResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		RenewAt:     fw.RenewAt,
		Deferred:    EphemeralResourceDeferred(fw.Deferred),
	}

	result, diags := EphemeralResultData(ctx, fw.Result)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.Result = result

	newPrivate, diags := fw.Private.Bytes(ctx)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.Private = newPrivate

	return proto6
}
