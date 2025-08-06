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

var _ ListValuable = &ListValue{}

// ListValuable extends attr.Value for list value types.
// Implement this interface to create a custom List value type.
type ListValuable interface {
	attr.Value

	// ToListValue should convert the value type to a List.
	ToListValue(ctx context.Context) (ListValue, diag.Diagnostics)
}

// ListValuableWithSemanticEquals extends ListValuable with semantic equality
// logic.
type ListValuableWithSemanticEquals interface {
	ListValuable

	// ListSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as computed elements added by
	// a remote system.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	ListSemanticEquals(context.Context, ListValuable) (bool, diag.Diagnostics)
}

// NewListNull creates a List with a null value. Determine whether the value is
// null via the List type IsNull method.
func NewListNull(elementType attr.Type) ListValue {
	return ListValue{
		elementType: elementType,
		state:       attr.ValueStateNull,
	}
}

// NewListUnknown creates a List with an unknown value. Determine whether the
// value is unknown via the List type IsUnknown method.
func NewListUnknown(elementType attr.Type) ListValue {
	return ListValue{
		elementType: elementType,
		state:       attr.ValueStateUnknown,
	}
}

// NewListValue creates a List with a known value. Access the value via the List
// type Elements or ElementsAs methods.
func NewListValue(elementType attr.Type, elements []attr.Value) (ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for idx, element := range elements {
		if !elementType.Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid List Element Type",
				"While creating a List value, an invalid element was detected. "+
					"A List must use the single, given element type. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("List Element Type: %s\n", elementType.String())+
					fmt.Sprintf("List Index (%d) Element Type: %s", idx, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewListUnknown(elementType), diags
	}

	return ListValue{
		elementType: elementType,
		elements:    elements,
		state:       attr.ValueStateKnown,
	}, nil
}

// NewListValueFrom creates a List with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the List type Elements or ElementsAs methods.
func NewListValueFrom(ctx context.Context, elementType attr.Type, elements any) (ListValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		ListType{ElemType: elementType},
		elements,
		path.Empty(),
	)

	if diags.HasError() {
		return NewListUnknown(elementType), diags
	}

	list, ok := attrValue.(ListValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert List Value",
			"An unexpected result occurred when creating a List using NewListValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return list, diags
}

// NewListValueMust creates a List with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the List
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create List values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewListValueMust(elementType attr.Type, elements []attr.Value) ListValue {
	list, diags := NewListValue(elementType, elements)

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

		panic("NewListValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return list
}

// ListValue represents a list of attr.Values, all of the same type, indicated
// by ElemType.
type ListValue struct {
	// elements is the collection of known values in the List.
	elements []attr.Value

	// elementType is the type of the elements in the List.
	elementType attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the collection of elements for the List.
func (l ListValue) Elements() []attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make([]attr.Value, 0, len(l.elements))
	result = append(result, l.elements...)

	return result
}

// ElementsAs populates `target` with the elements of the ListValue, throwing an
// error if the elements cannot be stored in `target`.
func (l ListValue) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this List to be able to use it with our
	// reflection code
	values, err := l.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"List Element Conversion Error",
				"An unexpected error was encountered trying to convert list elements. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}
	return reflect.Into(ctx, ListType{ElemType: l.elementType}, values, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// ElementType returns the element type for the List.
func (l ListValue) ElementType(_ context.Context) attr.Type {
	return l.elementType
}

// Type returns a ListType with the same element type as `l`.
func (l ListValue) Type(ctx context.Context) attr.Type {
	return ListType{ElemType: l.ElementType(ctx)}
}

// ToTerraformValue returns the data contained in the List as a tftypes.Value.
func (l ListValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	listType := tftypes.List{ElementType: l.ElementType(ctx).TerraformType(ctx)}

	switch l.state {
	case attr.ValueStateKnown:
		// MAINTAINER NOTE:
		// ListValue does not support DynamicType as an element type. It is not explicitly prevented from being created with the
		// Framework type system, but the Framework-supported ListAttribute, ListNestedAttribute, and ListNestedBlock all prevent DynamicType
		// from being used as an element type.
		//
		// In the future, if we ever need to support a list of dynamic element types, this tftypes.List creation logic will need to be modified to ensure
		// that known values contain the exact same concrete element type, specifically with unknown and null values. Dynamic values will return the correct concrete
		// element type for known values from `elem.ToTerraformValue`, but unknown and null values will be tftypes.DynamicPseudoType, causing an error due to multiple element
		// types in a tftypes.List.
		//
		// Unknown and null element types of tftypes.DynamicPseudoType must be recreated as the concrete element type unknown/null value. This can be done by checking `l.elements`
		// for a single concrete type (i.e. not tftypes.DynamicPseudoType), and using that concrete type to create unknown and null dynamic values later.
		//
		vals := make([]tftypes.Value, 0, len(l.elements))

		for _, elem := range l.elements {
			val, err := elem.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(listType, tftypes.UnknownValue), err
			}

			vals = append(vals, val)
		}

		if err := tftypes.ValidateValue(listType, vals); err != nil {
			return tftypes.NewValue(listType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(listType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(listType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(listType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled List state in ToTerraformValue: %s", l.state))
	}
}

// Equal returns true if the given attr.Value is also a ListValue, has the
// same element type, same value state, and contains exactly the element values
// as defined by the Equal method of the element type.
func (l ListValue) Equal(o attr.Value) bool {
	other, ok := o.(ListValue)

	if !ok {
		return false
	}

	// A list with no elementType is an invalid state
	if l.elementType == nil || other.elementType == nil {
		return false
	}

	if !l.elementType.Equal(other.elementType) {
		return false
	}

	if l.state != other.state {
		return false
	}

	if l.state != attr.ValueStateKnown {
		return true
	}

	if len(l.elements) != len(other.elements) {
		return false
	}

	for idx, lElem := range l.elements {
		otherElem := other.elements[idx]

		if !lElem.Equal(otherElem) {
			return false
		}
	}

	return true
}

// IsNull returns true if the List represents a null value.
func (l ListValue) IsNull() bool {
	return l.state == attr.ValueStateNull
}

// IsUnknown returns true if the List represents a currently unknown value.
// Returns false if the List has a known number of elements, even if all are
// unknown values.
func (l ListValue) IsUnknown() bool {
	return l.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the List value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (l ListValue) String() string {
	if l.IsUnknown() {
		return attr.UnknownValueString
	}

	if l.IsNull() {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range l.Elements() {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}

// ToListValue returns the List.
func (l ListValue) ToListValue(context.Context) (ListValue, diag.Diagnostics) {
	return l, nil
}
