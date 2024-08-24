// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueSemanticEqualitySet performs set type semantic equality.
//
// This will perform semantic equality checking on elements, regardless of
// whether the collection type implements the expected interface, since it
// cannot be assumed that the collection type implementation runs all possible
// element implementations.
func ValueSemanticEqualitySet(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.SetValuableWithSemanticEquals)

	// While the collection type itself does not implement the interface,
	// underlying elements might. Check elements automatically, if possible.
	if !ok {
		ValueSemanticEqualitySetElements(ctx, req, resp)

		return
	}

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.SetValuableWithSemanticEquals)

	// While the collection type itself does not implement the interface,
	// underlying elements might. Check elements automatically, if possible.
	if !ok {
		ValueSemanticEqualitySetElements(ctx, req, resp)

		return
	}

	logging.FrameworkTrace(
		ctx,
		"Calling provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	usePriorValue, diags := proposedNewValuable.SetSemanticEquals(ctx, priorValuable)

	logging.FrameworkTrace(
		ctx,
		"Called provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If the collection type signaled semantic equality, respect the
	// determination to use the whole prior value and return early since
	// checking elements is not necessary.
	if usePriorValue {
		resp.NewValue = priorValuable

		return
	}

	// While the collection type itself did not signal semantic equality,
	// underlying elements might, which should still modify the collection.
	// Check elements automatically, if possible.
	//
	// This logic pessimistically assumes that collection type semantic equality
	// implementations may be missing proper element type handling. While
	// correct implementations receive a small performance penalty of
	// being re-checked, this ensures that less-correct implementations do not
	// cause inconsistent data handling behaviors for developers.
	ValueSemanticEqualitySetElements(ctx, req, resp)
}

// ValueSemanticEqualitySetElements performs list type semantic equality
// on elements, returning a modified list as necessary.
func ValueSemanticEqualitySetElements(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.SetValuable)

	// No changes required if the elements cannot be extracted.
	if !ok {
		return
	}

	priorValue, diags := priorValuable.ToSetValue(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	priorValueElements := priorValue.Elements()

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.SetValuable)

	// No changes required if the elements cannot be extracted.
	if !ok {
		return
	}

	proposedNewValue, diags := proposedNewValuable.ToSetValue(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	proposedNewValueElements := proposedNewValue.Elements()

	// Create a new element value slice, which will be used to create the final
	// collection value after each element is evaluated.
	newValueElements := make([]attr.Value, len(proposedNewValueElements))

	// Short circuit flag
	updatedElements := false

	// Loop through proposed elements by delegating to the recursive semantic
	// equality logic. This ensures that recursion will catch a further
	// underlying element type has its semantic equality logic checked, even if
	// the current element type does not implement the interface.
	for idx, proposedNewValueElement := range proposedNewValueElements {
		// Ensure new value always contains all of proposed new value
		newValueElements[idx] = proposedNewValueElement

		if idx >= len(priorValueElements) {
			continue
		}

		elementReq := ValueSemanticEqualityRequest{
			Path:             req.Path.AtSetValue(proposedNewValueElement),
			PriorValue:       priorValueElements[idx],
			ProposedNewValue: proposedNewValueElement,
		}
		elementResp := &ValueSemanticEqualityResponse{
			NewValue: elementReq.ProposedNewValue,
		}

		ValueSemanticEquality(ctx, elementReq, elementResp)

		resp.Diagnostics.Append(elementResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}

		if elementResp.NewValue.Equal(elementReq.ProposedNewValue) {
			continue
		}

		updatedElements = true
		newValueElements[idx] = elementResp.NewValue
	}

	// No changes required if the elements were not updated.
	if !updatedElements {
		return
	}

	newValue, diags := basetypes.NewSetValue(proposedNewValue.ElementType(ctx), newValueElements)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the new value to the original SetValuable type to ensure
	// downstream logic has the correct value type for the defined schema type.
	newTypable, ok := proposedNewValuable.Type(ctx).(basetypes.SetTypable)

	// This should be a requirement of having a SetValuable, but defensively
	// checking just in case.
	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value Semantic Equality Type Error",
			"An unexpected error occurred while performing value semantic equality logic. "+
				"This is either an error in terraform-plugin-framework or a provider custom type implementation. "+
				"Please report this to the provider developers.\n\n"+
				"Error: Expected basetypes.SetTypable type for value type: "+fmt.Sprintf("%T", proposedNewValuable)+"\n"+
				"Path: "+req.Path.String(),
		)

		return
	}

	newValuable, diags := newTypable.ValueFromSet(ctx, newValue)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.NewValue = newValuable
}
