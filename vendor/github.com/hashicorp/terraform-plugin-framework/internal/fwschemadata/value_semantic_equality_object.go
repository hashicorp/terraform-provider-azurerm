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

// ValueSemanticEqualityObject performs object type semantic equality.
//
// This will perform semantic equality checking on attributes, regardless of
// whether the structural type implements the expected interface, since it
// cannot be assumed that the structural type implementation runs all possible
// attribute implementations.
func ValueSemanticEqualityObject(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.ObjectValuableWithSemanticEquals)

	// While the structural type itself does not implement the interface,
	// underlying attributes might. Check attributes automatically, if possible.
	if !ok {
		ValueSemanticEqualityObjectAttributes(ctx, req, resp)

		return
	}

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.ObjectValuableWithSemanticEquals)

	// While the structural type itself does not implement the interface,
	// underlying attributes might. Check attributes automatically, if possible.
	if !ok {
		ValueSemanticEqualityObjectAttributes(ctx, req, resp)

		return
	}

	logging.FrameworkTrace(
		ctx,
		"Calling provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	usePriorValue, diags := proposedNewValuable.ObjectSemanticEquals(ctx, priorValuable)

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

	// If the structural type signaled semantic equality, respect the
	// determination to use the whole prior value and return early since
	// checking attributes is not necessary.
	if usePriorValue {
		resp.NewValue = priorValuable

		return
	}

	// While the structural type itself did not signal semantic equality,
	// underlying attributes might, which should still modify the structural.
	// Check attributes automatically, if possible.
	//
	// This logic pessimistically assumes that structural type semantic equality
	// implementations may be missing proper attribute type handling. While
	// correct implementations receive a small performance penalty of
	// being re-checked, this ensures that less-correct implementations do not
	// cause inconsistent data handling behaviors for developers.
	ValueSemanticEqualityObjectAttributes(ctx, req, resp)
}

// ValueSemanticEqualityObjectAttributes performs object type semantic equality
// on attributes, returning a modified object as necessary.
func ValueSemanticEqualityObjectAttributes(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.ObjectValuable)

	// No changes required if the attributes cannot be extracted.
	if !ok {
		return
	}

	priorValue, diags := priorValuable.ToObjectValue(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	priorValueAttributes := priorValue.Attributes()

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.ObjectValuable)

	// No changes required if the attributes cannot be extracted.
	if !ok {
		return
	}

	proposedNewValue, diags := proposedNewValuable.ToObjectValue(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	proposedNewValueAttributes := proposedNewValue.Attributes()

	// Create a new element value map, which will be used to create the final
	// collection value after each element is evaluated.
	newValueAttributes := make(map[string]attr.Value, len(proposedNewValueAttributes))

	// Short circuit flag
	updatedAttributes := false

	// Loop through proposed attributes by delegating to the recursive semantic
	// equality logic. This ensures that recursion will catch a further
	// underlying element type has its semantic equality logic checked, even if
	// the current element type does not implement the interface.
	for name, proposedNewValueElement := range proposedNewValueAttributes {
		// Ensure new value always contains all of proposed new value
		newValueAttributes[name] = proposedNewValueElement

		priorValueElement, ok := priorValueAttributes[name]

		if !ok {
			continue
		}

		elementReq := ValueSemanticEqualityRequest{
			Path:             req.Path.AtName(name),
			PriorValue:       priorValueElement,
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

		updatedAttributes = true
		newValueAttributes[name] = elementResp.NewValue
	}

	// No changes required if the attributes were not updated.
	if !updatedAttributes {
		return
	}

	newValue, diags := basetypes.NewObjectValue(proposedNewValue.AttributeTypes(ctx), newValueAttributes)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the new value to the original ObjectValuable type to ensure
	// downstream logic has the correct value type for the defined schema type.
	newTypable, ok := proposedNewValuable.Type(ctx).(basetypes.ObjectTypable)

	// This should be a requirement of having a ObjectValuable, but defensively
	// checking just in case.
	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value Semantic Equality Type Error",
			"An unexpected error occurred while performing value semantic equality logic. "+
				"This is either an error in terraform-plugin-framework or a provider custom type implementation. "+
				"Please report this to the provider developers.\n\n"+
				"Error: Expected basetypes.ObjectTypable type for value type: "+fmt.Sprintf("%T", proposedNewValuable)+"\n"+
				"Path: "+req.Path.String(),
		)

		return
	}

	newValuable, diags := newTypable.ValueFromObject(ctx, newValue)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.NewValue = newValuable
}
