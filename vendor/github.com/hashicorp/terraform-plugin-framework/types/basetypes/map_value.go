// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ MapValuable = &MapValue{}

// MapValuable extends attr.Value for map value types.
// Implement this interface to create a custom Map value type.
type MapValuable interface {
	attr.Value

	// ToMapValue should convert the value type to a Map.
	ToMapValue(ctx context.Context) (MapValue, diag.Diagnostics)
}

// MapValuableWithSemanticEquals extends MapValuable with semantic equality
// logic.
type MapValuableWithSemanticEquals interface {
	MapValuable

	// MapSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as computed elements added by
	// a remote system.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	MapSemanticEquals(context.Context, MapValuable) (bool, diag.Diagnostics)
}

// NewMapNull creates a Map with a null value. Determine whether the value is
// null via the Map type IsNull method.
func NewMapNull(elementType attr.Type) MapValue {
	return MapValue{
		elementType: elementType,
		state:       attr.ValueStateNull,
	}
}

// NewMapUnknown creates a Map with an unknown value. Determine whether the
// value is unknown via the Map type IsUnknown method.
func NewMapUnknown(elementType attr.Type) MapValue {
	return MapValue{
		elementType: elementType,
		state:       attr.ValueStateUnknown,
	}
}

// NewMapValue creates a Map with a known value. Access the value via the Map
// type Elements or ElementsAs methods.
func NewMapValue(elementType attr.Type, elements map[string]attr.Value) (MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for key, element := range elements {
		if !elementType.Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid Map Element Type",
				"While creating a Map value, an invalid element was detected. "+
					"A Map must use the single, given element type. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Map Element Type: %s\n", elementType.String())+
					fmt.Sprintf("Map Key (%s) Element Type: %s", key, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewMapUnknown(elementType), diags
	}

	return MapValue{
		elementType: elementType,
		elements:    elements,
		state:       attr.ValueStateKnown,
	}, nil
}

// NewMapValueFrom creates a Map with a known value, using reflection rules.
// The elements must be a map of string keys to values which can convert into
// the given element type. Access the value via the Map type Elements or
// ElementsAs methods.
func NewMapValueFrom(ctx context.Context, elementType attr.Type, elements any) (MapValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		MapType{ElemType: elementType},
		elements,
		path.Empty(),
	)

	if diags.HasError() {
		return NewMapUnknown(elementType), diags
	}

	m, ok := attrValue.(MapValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert Map Value",
			"An unexpected result occurred when creating a Map using MapValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return m, diags
}

// NewMapValueMust creates a Map with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Map
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Map values which will
// not potentially effect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewMapValueMust(elementType attr.Type, elements map[string]attr.Value) MapValue {
	m, diags := NewMapValue(elementType, elements)

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

		panic("MapValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return m
}

// MapValue represents a mapping of string keys to attr.Value values of a single
// type.
type MapValue struct {
	// elements is the mapping of known values in the Map.
	elements map[string]attr.Value

	// elementType is the type of the elements in the Map.
	elementType attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the mapping of elements for the Map.
func (m MapValue) Elements() map[string]attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make(map[string]attr.Value, len(m.elements))

	for key, value := range m.elements {
		result[key] = value
	}

	return result
}

// ElementsAs populates `target` with the elements of the MapValue, throwing an
// error if the elements cannot be stored in `target`.
func (m MapValue) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this Map to be able to use it with our
	// reflection code
	val, err := m.ToTerraformValue(ctx)
	if err != nil {
		err := fmt.Errorf("error getting Terraform value for map: %w", err)
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Map Conversion Error",
				"An unexpected error was encountered trying to convert the map into an equivalent Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}

	return reflect.Into(ctx, MapType{ElemType: m.elementType}, val, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// ElementType returns the element type for the Map.
func (m MapValue) ElementType(_ context.Context) attr.Type {
	return m.elementType
}

// Type returns a MapType with the same element type as `m`.
func (m MapValue) Type(ctx context.Context) attr.Type {
	return MapType{ElemType: m.ElementType(ctx)}
}

// ToTerraformValue returns the data contained in the Map as a tftypes.Value.
func (m MapValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	mapType := tftypes.Map{ElementType: m.ElementType(ctx).TerraformType(ctx)}

	switch m.state {
	case attr.ValueStateKnown:
		// MAINTAINER NOTE:
		// MapValue does not support DynamicType as an element type. It is not explicitly prevented from being created with the
		// Framework type system, but the Framework-supported MapAttribute and MapNestedAttribute prevent DynamicType
		// from being used as an element type.
		//
		// In the future, if we ever need to support a map of dynamic element types, this tftypes.Map creation logic will need to be modified to ensure
		// that known values contain the exact same concrete element type, specifically with unknown and null values. Dynamic values will return the correct concrete
		// element type for known values from `elem.ToTerraformValue`, but unknown and null values will be tftypes.DynamicPseudoType, causing an error due to multiple element
		// types in a tftypes.Map.
		//
		// Unknown and null element types of tftypes.DynamicPseudoType must be recreated as the concrete element type unknown/null value. This can be done by checking `m.elements`
		// for a single concrete type (i.e. not tftypes.DynamicPseudoType), and using that concrete type to create unknown and null dynamic values later.
		//
		vals := make(map[string]tftypes.Value, len(m.elements))

		for key, elem := range m.elements {
			val, err := elem.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(mapType, tftypes.UnknownValue), err
			}

			vals[key] = val
		}

		if err := tftypes.ValidateValue(mapType, vals); err != nil {
			return tftypes.NewValue(mapType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(mapType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(mapType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(mapType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Map state in ToTerraformValue: %s", m.state))
	}
}

// Equal returns true if the given attr.Value is also a MapValue, has the
// same element type, same value state, and contains exactly the element values
// as defined by the Equal method of the element type.
func (m MapValue) Equal(o attr.Value) bool {
	other, ok := o.(MapValue)

	if !ok {
		return false
	}

	// A map with no elementType is an invalid state
	if m.elementType == nil || other.elementType == nil {
		return false
	}

	if !m.elementType.Equal(other.elementType) {
		return false
	}

	if m.state != other.state {
		return false
	}

	if m.state != attr.ValueStateKnown {
		return true
	}

	if len(m.elements) != len(other.elements) {
		return false
	}

	for key, mElem := range m.elements {
		otherElem := other.elements[key]

		if !mElem.Equal(otherElem) {
			return false
		}
	}

	return true
}

// IsNull returns true if the Map represents a null value.
func (m MapValue) IsNull() bool {
	return m.state == attr.ValueStateNull
}

// IsUnknown returns true if the Map represents a currently unknown value.
// Returns false if the Map has a known number of elements, even if all are
// unknown values.
func (m MapValue) IsUnknown() bool {
	return m.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Map value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (m MapValue) String() string {
	if m.IsUnknown() {
		return attr.UnknownValueString
	}

	if m.IsNull() {
		return attr.NullValueString
	}

	// We want the output to be consistent, so we sort the output by key
	keys := make([]string, 0, len(m.Elements()))
	for k := range m.Elements() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res strings.Builder

	res.WriteString("{")
	for i, k := range keys {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(fmt.Sprintf("%q:%s", k, m.Elements()[k].String()))
	}
	res.WriteString("}")

	return res.String()
}

// ToMapValue returns the Map.
func (m MapValue) ToMapValue(context.Context) (MapValue, diag.Diagnostics) {
	return m, nil
}
