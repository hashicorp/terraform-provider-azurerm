// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ SetValuable = &SetValue{}

// SetValuable extends attr.Value for set value types.
// Implement this interface to create a custom Set value type.
type SetValuable interface {
	attr.Value

	// ToSetValue should convert the value type to a Set.
	ToSetValue(ctx context.Context) (SetValue, diag.Diagnostics)
}

// SetValuableWithSemanticEquals extends SetValuable with semantic equality
// logic.
type SetValuableWithSemanticEquals interface {
	SetValuable

	// SetSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as computed elements added by
	// a remote system.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	SetSemanticEquals(context.Context, SetValuable) (bool, diag.Diagnostics)
}

// NewSetNull creates a Set with a null value. Determine whether the value is
// null via the Set type IsNull method.
func NewSetNull(elementType attr.Type) SetValue {
	return SetValue{
		elementType: elementType,
		state:       attr.ValueStateNull,
	}
}

// NewSetUnknown creates a Set with an unknown value. Determine whether the
// value is unknown via the Set type IsUnknown method.
func NewSetUnknown(elementType attr.Type) SetValue {
	return SetValue{
		elementType: elementType,
		state:       attr.ValueStateUnknown,
	}
}

// NewSetValue creates a Set with a known value. Access the value via the Set
// type Elements or ElementsAs methods.
func NewSetValue(elementType attr.Type, elements []attr.Value) (SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for idx, element := range elements {
		if !elementType.Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid Set Element Type",
				"While creating a Set value, an invalid element was detected. "+
					"A Set must use the single, given element type. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Set Element Type: %s\n", elementType.String())+
					fmt.Sprintf("Set Index (%d) Element Type: %s", idx, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewSetUnknown(elementType), diags
	}

	return SetValue{
		elementType: elementType,
		elements:    elements,
		state:       attr.ValueStateKnown,
	}, nil
}

// NewSetValueFrom creates a Set with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the Set type Elements or ElementsAs methods.
func NewSetValueFrom(ctx context.Context, elementType attr.Type, elements any) (SetValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		SetType{ElemType: elementType},
		elements,
		path.Empty(),
	)

	if diags.HasError() {
		return NewSetUnknown(elementType), diags
	}

	set, ok := attrValue.(SetValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert Set Value",
			"An unexpected result occurred when creating a Set using SetValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return set, diags
}

// NewSetValueMust creates a Set with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Set
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Set values which will
// not potentially effect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewSetValueMust(elementType attr.Type, elements []attr.Value) SetValue {
	set, diags := NewSetValue(elementType, elements)

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

		panic("SetValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return set
}

// SetValue represents a set of attr.Value, all of the same type,
// indicated by ElemType.
type SetValue struct {
	// elements is the collection of known values in the Set.
	elements []attr.Value

	// elementType is the type of the elements in the Set.
	elementType attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the collection of elements for the Set.
func (s SetValue) Elements() []attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make([]attr.Value, 0, len(s.elements))
	result = append(result, s.elements...)

	return result
}

// ElementsAs populates `target` with the elements of the SetValue, throwing an
// error if the elements cannot be stored in `target`.
func (s SetValue) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this Set to be able to use it with our
	// reflection code
	val, err := s.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Set Element Conversion Error",
				"An unexpected error was encountered trying to convert set elements. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}
	return reflect.Into(ctx, s.Type(ctx), val, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// ElementType returns the element type for the Set.
func (s SetValue) ElementType(_ context.Context) attr.Type {
	return s.elementType
}

// Type returns a SetType with the same element type as `s`.
func (s SetValue) Type(ctx context.Context) attr.Type {
	return SetType{ElemType: s.ElementType(ctx)}
}

// ToTerraformValue returns the data contained in the Set as a tftypes.Value.
func (s SetValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	setType := tftypes.Set{ElementType: s.ElementType(ctx).TerraformType(ctx)}

	switch s.state {
	case attr.ValueStateKnown:
		// MAINTAINER NOTE:
		// SetValue does not support DynamicType as an element type. It is not explicitly prevented from being created with the
		// Framework type system, but the Framework-supported SetAttribute, SetNestedAttribute, and SetNestedBlock all prevent DynamicType
		// from being used as an element type.
		//
		// In the future, if we ever need to support a set of dynamic element types, this tftypes.Set creation logic will need to be modified to ensure
		// that known values contain the exact same concrete element type, specifically with unknown and null values. Dynamic values will return the correct concrete
		// element type for known values from `elem.ToTerraformValue`, but unknown and null values will be tftypes.DynamicPseudoType, causing an error due to multiple element
		// types in a tftypes.Set.
		//
		// Unknown and null element types of tftypes.DynamicPseudoType must be recreated as the concrete element type unknown/null value. This can be done by checking `s.elements`
		// for a single concrete type (i.e. not tftypes.DynamicPseudoType), and using that concrete type to create unknown and null dynamic values later.
		//
		vals := make([]tftypes.Value, 0, len(s.elements))

		for _, elem := range s.elements {
			val, err := elem.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(setType, tftypes.UnknownValue), err
			}

			vals = append(vals, val)
		}

		if err := tftypes.ValidateValue(setType, vals); err != nil {
			return tftypes.NewValue(setType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(setType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(setType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(setType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Set state in ToTerraformValue: %s", s.state))
	}
}

// Equal returns true if the given attr.Value is also a SetValue, has the
// same element type, same value state, and contains exactly the element values
// as defined by the Equal method of the element type.
func (s SetValue) Equal(o attr.Value) bool {
	other, ok := o.(SetValue)

	if !ok {
		return false
	}

	// A set with no elementType is an invalid state
	if s.elementType == nil || other.elementType == nil {
		return false
	}

	if !s.elementType.Equal(other.elementType) {
		return false
	}

	if s.state != other.state {
		return false
	}

	if s.state != attr.ValueStateKnown {
		return true
	}

	if len(s.elements) != len(other.elements) {
		return false
	}

	for _, elem := range s.elements {
		if !other.contains(elem) {
			return false
		}
	}

	return true
}

func (s SetValue) contains(v attr.Value) bool {
	for _, elem := range s.Elements() {
		if elem.Equal(v) {
			return true
		}
	}

	return false
}

// IsNull returns true if the Set represents a null value.
func (s SetValue) IsNull() bool {
	return s.state == attr.ValueStateNull
}

// IsUnknown returns true if the Set represents a currently unknown value.
// Returns false if the Set has a known number of elements, even if all are
// unknown values.
func (s SetValue) IsUnknown() bool {
	return s.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Set value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (s SetValue) String() string {
	if s.IsUnknown() {
		return attr.UnknownValueString
	}

	if s.IsNull() {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range s.Elements() {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}

// ToSetValue returns the Set.
func (s SetValue) ToSetValue(context.Context) (SetValue, diag.Diagnostics) {
	return s, nil
}
