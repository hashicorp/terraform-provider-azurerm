// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ attr.Value = TupleValue{}

// NewTupleNull creates a Tuple with a null value.
func NewTupleNull(elementTypes []attr.Type) TupleValue {
	return TupleValue{
		elementTypes: elementTypes,
		state:        attr.ValueStateNull,
	}
}

// NewTupleUnknown creates a Tuple with an unknown value.
func NewTupleUnknown(elementTypes []attr.Type) TupleValue {
	return TupleValue{
		elementTypes: elementTypes,
		state:        attr.ValueStateUnknown,
	}
}

// NewTupleValue creates a Tuple with a known value. Access the value via the Tuple type Elements method.
func NewTupleValue(elementTypes []attr.Type, elements []attr.Value) (TupleValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	if len(elementTypes) != len(elements) {
		givenTypes := make([]attr.Type, len(elements))
		for i, v := range elements {
			givenTypes[i] = v.Type(ctx)
		}

		diags.AddError(
			"Invalid Tuple Elements",
			"While creating a Tuple value, mismatched element types were detected. "+
				"A Tuple must be an ordered array of elements where the values exactly match the length and types of the defined element types. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Tuple Expected Type: %v\n", elementTypes)+
				fmt.Sprintf("Tuple Given Type: %v", givenTypes),
		)

		return NewTupleUnknown(elementTypes), diags
	}

	for i, element := range elements {
		if !elementTypes[i].Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid Tuple Element",
				"While creating a Tuple value, an invalid element was detected. "+
					"A Tuple must be an ordered array of elements where the values exactly match the length and types of the defined element types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Tuple Index (%d) Expected Type: %s\n", i, elementTypes[i])+
					fmt.Sprintf("Tuple Index (%d) Given Type: %s", i, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewTupleUnknown(elementTypes), diags
	}

	return TupleValue{
		elementTypes: elementTypes,
		elements:     elements,
		state:        attr.ValueStateKnown,
	}, nil
}

// NewTupleValueMust creates a Tuple with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Tuple type Elements method.
//
// This creation function is only recommended to create Tuple values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewTupleValueMust(elementTypes []attr.Type, elements []attr.Value) TupleValue {
	tuple, diags := NewTupleValue(elementTypes, elements)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTupleValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return tuple
}

// TupleValue represents an ordered list of attr.Value, with an attr.Type for each element. This type intentionally
// includes less functionality than other types in the type system as it has limited real world application and therefore
// is not exposed to provider developers.
type TupleValue struct {
	// elements is the ordered list of known element values for the tuple.
	elements []attr.Value

	// elementTypes is the ordered list of elements types for the tuple.
	elementTypes []attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the ordered list of known values for the Tuple.
func (v TupleValue) Elements() []attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make([]attr.Value, 0, len(v.elements))
	result = append(result, v.elements...)

	return result
}

// ElementTypes returns the ordered list of element types for the Tuple.
func (v TupleValue) ElementTypes(ctx context.Context) []attr.Type {
	return v.elementTypes
}

// Equal returns true if the given attr.Value is also a Tuple, has the same value state,
// and contains exactly the same element types/values as defined by the Equal method of those
// underlying types/values.
func (v TupleValue) Equal(o attr.Value) bool {
	other, ok := o.(TupleValue)
	if !ok {
		return false
	}

	if len(v.elementTypes) != len(other.elementTypes) {
		return false
	}

	for i, elementType := range v.elementTypes {
		if !elementType.Equal(other.elementTypes[i]) {
			return false
		}
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	// This statement should never be true, given that element type length must exactly match the number of elements,
	// but checking to avoid an index out of range panic
	if len(v.elements) != len(other.elements) {
		return false
	}

	for i, element := range v.elements {
		if !element.Equal(other.elements[i]) {
			return false
		}
	}

	return true
}

// IsNull returns true if the Tuple represents a null value.
func (v TupleValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

// IsUnknown returns true if the Tuple represents an unknown value.
func (v TupleValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Tuple. The string returned here is not protected by any
// compatibility guarantees, and is intended for logging and error reporting.
func (v TupleValue) String() string {
	if v.IsUnknown() {
		return attr.UnknownValueString
	}

	if v.IsNull() {
		return attr.NullValueString
	}

	elements := v.Elements()
	valueStrings := make([]string, len(elements))

	for i, element := range elements {
		valueStrings[i] = element.String()
	}

	return "[" + strings.Join(valueStrings, ",") + "]"
}

// Type returns a TupleType with the elements types for the Tuple.
func (v TupleValue) Type(ctx context.Context) attr.Type {
	return TupleType{
		ElemTypes: v.ElementTypes(ctx),
	}
}

// ToTerraformValue returns the equivalent tftypes.Value for the Tuple.
func (v TupleValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	tfTypes := make([]tftypes.Type, len(v.elementTypes))
	for i, elementType := range v.elementTypes {
		tfTypes[i] = elementType.TerraformType(ctx)
	}

	tupleType := tftypes.Tuple{ElementTypes: tfTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make([]tftypes.Value, 0, len(v.elements))

		for _, elem := range v.elements {
			val, err := elem.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(tupleType, tftypes.UnknownValue), err
			}

			vals = append(vals, val)
		}

		if err := tftypes.ValidateValue(tupleType, vals); err != nil {
			return tftypes.NewValue(tupleType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tupleType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tupleType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tupleType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Tuple state in ToTerraformValue: %s", v.state))
	}
}
