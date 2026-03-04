// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ListResourceResult(ctx context.Context, result *fwserver.ListResult) tfprotov6.ListResourceResult {
	allDiags := result.Diagnostics
	if allDiags.HasError() {
		return tfprotov6.ListResourceResult{
			Diagnostics: Diagnostics(ctx, allDiags),
		}
	}

	resourceIdentity, diags := ResourceIdentity(ctx, result.Identity)
	allDiags.Append(diags...)

	return tfprotov6.ListResourceResult{
		DisplayName: result.DisplayName,
		Identity:    resourceIdentity,
		Diagnostics: Diagnostics(ctx, allDiags),
	}
}

func ListResourceResultWithResource(ctx context.Context, result *fwserver.ListResult) tfprotov6.ListResourceResult {
	allDiags := result.Diagnostics
	if allDiags.HasError() {
		return tfprotov6.ListResourceResult{
			Diagnostics: Diagnostics(ctx, allDiags),
		}
	}

	resourceIdentity, diags := ResourceIdentity(ctx, result.Identity)
	allDiags.Append(diags...)

	resource, diags := Resource(ctx, result.Resource)
	allDiags.Append(diags...)

	return tfprotov6.ListResourceResult{
		DisplayName: result.DisplayName,
		Identity:    resourceIdentity,
		Resource:    resource,
		Diagnostics: Diagnostics(ctx, allDiags),
	}
}
