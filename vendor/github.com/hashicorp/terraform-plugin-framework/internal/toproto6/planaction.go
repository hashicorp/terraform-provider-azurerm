// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// PlanActionResponse returns the *tfprotov6.PlanActionResponse equivalent of a *fwserver.PlanActionResponse.
func PlanActionResponse(ctx context.Context, fw *fwserver.PlanActionResponse) *tfprotov6.PlanActionResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.PlanActionResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		Deferred:    ActionDeferred(fw.Deferred),
	}

	return proto6
}
