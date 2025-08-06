// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueSemanticEqualityDynamic performs dynamic type semantic equality.
func ValueSemanticEqualityDynamic(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.DynamicValuableWithSemanticEquals)

	// No changes required if the interface is not implemented.
	if !ok {
		return
	}

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.DynamicValuableWithSemanticEquals)

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

	// The prior dynamic value has alredy been checked for null or unknown, however, we also
	// need to check the underlying value for null or unknown.
	priorValue, diags := priorValuable.ToDynamicValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if priorValue.IsUnderlyingValueNull() || priorValue.IsUnderlyingValueUnknown() {
		return
	}

	// The proposed new dynamic value has alredy been checked for null or unknown, however, we also
	// need to check the underlying value for null or unknown.
	proposedValue, diags := proposedNewValuable.ToDynamicValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if proposedValue.IsUnderlyingValueNull() || proposedValue.IsUnderlyingValueUnknown() {
		return
	}

	usePriorValue, diags := proposedNewValuable.DynamicSemanticEquals(ctx, priorValuable)

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
