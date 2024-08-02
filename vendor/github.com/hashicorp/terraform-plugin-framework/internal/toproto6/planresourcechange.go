// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
)

// PlanResourceChangeResponse returns the *tfprotov6.PlanResourceChangeResponse
// equivalent of a *fwserver.PlanResourceChangeResponse.
func PlanResourceChangeResponse(ctx context.Context, fw *fwserver.PlanResourceChangeResponse) *tfprotov6.PlanResourceChangeResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.PlanResourceChangeResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	plannedState, diags := State(ctx, fw.PlannedState)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.PlannedState = plannedState

	requiresReplace, diags := totftypes.AttributePaths(ctx, fw.RequiresReplace)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.RequiresReplace = requiresReplace

	plannedPrivate, diags := fw.PlannedPrivate.Bytes(ctx)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.PlannedPrivate = plannedPrivate

	return proto6
}
