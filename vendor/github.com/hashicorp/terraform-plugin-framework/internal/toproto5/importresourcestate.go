// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ImportResourceStateResponse returns the *tfprotov5.ImportResourceStateResponse
// equivalent of a *fwserver.ImportResourceStateResponse.
func ImportResourceStateResponse(ctx context.Context, fw *fwserver.ImportResourceStateResponse) *tfprotov5.ImportResourceStateResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ImportResourceStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	for _, fwImportedResource := range fw.ImportedResources {
		proto5ImportedResource, diags := ImportedResource(ctx, &fwImportedResource)

		proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)

		if diags.HasError() {
			continue
		}

		proto5.ImportedResources = append(proto5.ImportedResources, proto5ImportedResource)
	}

	return proto5
}
