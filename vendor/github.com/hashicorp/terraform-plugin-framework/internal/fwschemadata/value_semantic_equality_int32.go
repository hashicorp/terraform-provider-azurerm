// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueSemanticEqualityInt32 performs int32 type semantic equality.
func ValueSemanticEqualityInt32(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.Int32ValuableWithSemanticEquals)

	// No changes required if the interface is not implemented.
	if !ok {
		return
	}

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.Int32ValuableWithSemanticEquals)

	// No changes required if the interface is not implemented.
	if !ok {
		return
	}

	logging.FrameworkTrace(
		ctx,
		"Calling provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	usePriorValue, diags := proposedNewValuable.Int32SemanticEquals(ctx, priorValuable)

	logging.FrameworkTrace(
		ctx,
		"Called provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	resp.Diagnostics.Append(diags...)

	if !usePriorValue {
		return
	}

	resp.NewValue = priorValuable
}
