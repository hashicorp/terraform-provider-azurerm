// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// CreateParentTerraformValue ensures that the given parent value can have children
// values upserted. If the parent value is known and not null, it is returned
// without modification. A null Object or Tuple is converted to known with null
// children. An unknown Object or Tuple is converted to known with unknown
// children. List, Map, and Set are created with empty elements.
func CreateParentTerraformValue(_ context.Context, parentPath path.Path, parentType tftypes.Type, childValue interface{}) (tftypes.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	var parentValue tftypes.Value

	switch parentType := parentType.(type) {
	case tftypes.List:
		parentValue = tftypes.NewValue(parentType, []tftypes.Value{})
	case tftypes.Set:
		parentValue = tftypes.NewValue(parentType, []tftypes.Value{})
	case tftypes.Map:
		parentValue = tftypes.NewValue(parentType, map[string]tftypes.Value{})
	case tftypes.Object:
		vals := map[string]tftypes.Value{}

		for name, t := range parentType.AttributeTypes {
			vals[name] = tftypes.NewValue(t, childValue)
		}

		parentValue = tftypes.NewValue(parentType, vals)
	case tftypes.Tuple:
		vals := []tftypes.Value{}

		for _, elementType := range parentType.ElementTypes {
			vals = append(vals, tftypes.NewValue(elementType, childValue))
		}

		parentValue = tftypes.NewValue(parentType, vals)
	default:
		diags.AddAttributeError(
			parentPath,
			"Value Conversion Error",
			"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Unknown parent type %s to create value.", parentType),
		)
		return parentValue, diags
	}

	return parentValue, diags
}

// UpsertChildTerraformValue will upsert a child value into a parent value. If the
// path step already has a value, it will be overwritten. Otherwise, the child
// value will be added.
//
// Lists can only have the next element added according to the current length.
func UpsertChildTerraformValue(_ context.Context, parentPath path.Path, parentValue tftypes.Value, childStep path.PathStep, childValue tftypes.Value) (tftypes.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO: Add Tuple support
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/54
	switch childStep := childStep.(type) {
	case path.PathStepAttributeName:
		// Set in Object
		if !parentValue.Type().Is(tftypes.Object{}) {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Cannot add attribute into parent type: %s", parentValue.Type()),
			)
			return parentValue, diags
		}

		var parentAttrs map[string]tftypes.Value
		err := parentValue.Copy().As(&parentAttrs)

		if err != nil {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Unable to extract object elements from parent value: %s", err),
			)
			return parentValue, diags
		}

		parentAttrs[string(childStep)] = childValue
		parentValue = tftypes.NewValue(parentValue.Type(), parentAttrs)
	case path.PathStepElementKeyInt:
		// Upsert List element, except past length + 1
		if !parentValue.Type().Is(tftypes.List{}) {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Cannot add list element into parent type: %s", parentValue.Type()),
			)
			return parentValue, diags
		}

		var parentElems []tftypes.Value
		err := parentValue.Copy().As(&parentElems)

		if err != nil {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Unable to extract list elements from parent value: %s", err),
			)
			return parentValue, diags
		}

		if int(childStep) > len(parentElems) {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Cannot add list element %d as list currently has %d length. To prevent ambiguity, only the next element can be added to a list. Add empty elements into the list prior to this call, if appropriate.", int(childStep)+1, len(parentElems)),
			)
			return parentValue, diags
		}

		if int(childStep) == len(parentElems) {
			parentElems = append(parentElems, childValue)
		} else {
			parentElems[int(childStep)] = childValue
		}

		parentValue = tftypes.NewValue(parentValue.Type(), parentElems)
	case path.PathStepElementKeyString:
		// Upsert Map element
		if !parentValue.Type().Is(tftypes.Map{}) {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Cannot add map value into parent type: %s", parentValue.Type()),
			)
			return parentValue, diags
		}

		var parentElems map[string]tftypes.Value
		err := parentValue.Copy().As(&parentElems)

		if err != nil {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Unable to extract map elements from parent value: %s", err),
			)
			return parentValue, diags
		}

		parentElems[string(childStep)] = childValue
		parentValue = tftypes.NewValue(parentValue.Type(), parentElems)
	case path.PathStepElementKeyValue:
		// Upsert Set element
		if !parentValue.Type().Is(tftypes.Set{}) {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Cannot add set element into parent type: %s", parentValue.Type()),
			)
			return parentValue, diags
		}

		var parentElems []tftypes.Value
		err := parentValue.Copy().As(&parentElems)

		if err != nil {
			diags.AddAttributeError(
				parentPath,
				"Value Conversion Error",
				"An unexpected error was encountered trying to create a value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
					fmt.Sprintf("Unable to extract set elements from parent value: %s", err),
			)
			return parentValue, diags
		}

		// Prevent duplicates
		var found bool

		for _, parentElem := range parentElems {
			if parentElem.Equal(childValue) {
				found = true
				break
			}
		}

		if !found {
			parentElems = append(parentElems, childValue)
		}

		parentValue = tftypes.NewValue(parentValue.Type(), parentElems)
	}

	return parentValue, diags
}
