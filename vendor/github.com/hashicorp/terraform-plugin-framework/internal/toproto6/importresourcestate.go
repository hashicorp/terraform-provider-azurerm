// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ImportResourceStateResponse returns the *tfprotov6.ImportResourceStateResponse
// equivalent of a *fwserver.ImportResourceStateResponse.
func ImportResourceStateResponse(ctx context.Context, fw *fwserver.ImportResourceStateResponse) *tfprotov6.ImportResourceStateResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ImportResourceStateResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	for _, fwImportedResource := range fw.ImportedResources {
		proto6ImportedResource, diags := ImportedResource(ctx, &fwImportedResource)

		proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)

		if diags.HasError() {
			continue
		}

		proto6.ImportedResources = append(proto6.ImportedResources, proto6ImportedResource)
	}

	return proto6
}
