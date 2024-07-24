// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
)

// PlanResourceChangeResponse returns the *tfprotov5.PlanResourceChangeResponse
// equivalent of a *fwserver.PlanResourceChangeResponse.
func PlanResourceChangeResponse(ctx context.Context, fw *fwserver.PlanResourceChangeResponse) *tfprotov5.PlanResourceChangeResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.PlanResourceChangeResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	plannedState, diags := State(ctx, fw.PlannedState)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.PlannedState = plannedState

	requiresReplace, diags := totftypes.AttributePaths(ctx, fw.RequiresReplace)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.RequiresReplace = requiresReplace

	plannedPrivate, diags := fw.PlannedPrivate.Bytes(ctx)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.PlannedPrivate = plannedPrivate

	return proto5
}
