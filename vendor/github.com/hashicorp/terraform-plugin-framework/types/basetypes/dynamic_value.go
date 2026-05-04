// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ DynamicValuable = DynamicValue{}
)

// DynamicValuable extends attr.Value for dynamic value types. Implement this interface
// to create a custom Dynamic value type.
type DynamicValuable interface {
	attr.Value

	// ToDynamicValue should convert the value type to a DynamicValue.
	ToDynamicValue(context.Context) (DynamicValue, diag.Diagnostics)
}

// DynamicValuableWithSemanticEquals extends DynamicValuable with semantic equality logic.
type DynamicValuableWithSemanticEquals interface {
	DynamicValuable

	// DynamicSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	DynamicSemanticEquals(context.Context, DynamicValuable) (bool, diag.Diagnostics)
}

// NewDynamicValue creates a Dynamic with a known value. Access the value via the Dynamic
// type UnderlyingValue method. The concrete value type returned to Terraform from this value
// will be determined by the provided `(attr.Value).ToTerraformValue` function.
func NewDynamicValue(value attr.Value) DynamicValue {
	if value == nil {
		return NewDynamicNull()
	}

	return DynamicValue{
		value: value,
		state: attr.ValueStateKnown,
	}
}

// NewDynamicNull creates a Dynamic with a null value. The concrete value type returned to Terraform
// from this value will be tftypes.DynamicPseudoType.
func NewDynamicNull() DynamicValue {
	return DynamicValue{
		state: attr.ValueStateNull,
	}
}

// NewDynamicUnknown creates a Dynamic with an unknown value. The concrete value type returned to Terraform
// from this value will be tftypes.DynamicPseudoType.
func NewDynamicUnknown() DynamicValue {
	return DynamicValue{
		state: attr.ValueStateUnknown,
	}
}

// DynamicValue represents a dynamic value. Static types are always
// preferable over dynamic types in Terraform as practitioners will receive less
// helpful configuration assistance from validation error diagnostics and editor
// integrations.
type DynamicValue struct {
	// value contains the known value, if not null or unknown.
	value attr.Value

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Type returns DynamicType.
func (v DynamicValue) Type(ctx context.Context) attr.Type {
	return DynamicType{}
}

// ToTerraformValue returns the equivalent tftypes.Value for the DynamicValue.
func (v DynamicValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	switch v.state {
	case attr.ValueStateKnown:
		if v.value == nil {
			return tftypes.NewValue(tftypes.DynamicPseudoType, tftypes.UnknownValue),
				errors.New("invalid Dynamic state in ToTerraformValue: DynamicValue is known but the underlying value is unset")
		}

		return v.value.ToTerraformValue(ctx)
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.DynamicPseudoType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.DynamicPseudoType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Dynamic state in ToTerraformValue: %s", v.state))
	}
}

// Equal returns true if the given attr.Value is also a DynamicValue and contains an equal underlying value as defined by its Equal method.
func (v DynamicValue) Equal(o attr.Value) bool {
	other, ok := o.(DynamicValue)
	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	// Prevent panic and force inequality if either underlying value is nil
	if v.value == nil || other.value == nil {
		return false
	}

	return v.value.Equal(other.value)
}

// IsNull returns true if the DynamicValue represents a null value.
func (v DynamicValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

// IsUnknown returns true if the DynamicValue represents an unknown value.
func (v DynamicValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the DynamicValue. The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (v DynamicValue) String() string {
	if v.IsUnknown() {
		return attr.UnknownValueString
	}

	if v.IsNull() {
		return attr.NullValueString
	}

	if v.value == nil {
		return attr.UnsetValueString
	}

	return v.value.String()
}

// ToDynamicValue returns DynamicValue.
func (v DynamicValue) ToDynamicValue(ctx context.Context) (DynamicValue, diag.Diagnostics) {
	return v, nil
}

// UnderlyingValue returns the concrete underlying value in the DynamicValue. This will return `nil`
// if DynamicValue is null or unknown.
//
// A known DynamicValue can have an underlying value that is in null or unknown state in the
// scenario that the underlying value type has been refined by Terraform.
func (v DynamicValue) UnderlyingValue() attr.Value {
	return v.value
}

// IsUnderlyingValueNull is a helper method that return true only in the case where the underlying value has a
// known type but the value is null. This method will return false if the underlying type is not known.
//
// IsNull should be used to determine if the dynamic value does not have a known type and the value is null.
//
// An example of a known type with a null underlying value would be:
//
//	types.DynamicValue(types.StringNull())
func (v DynamicValue) IsUnderlyingValueNull() bool {
	return v.value != nil && v.value.IsNull()
}

// IsUnderlyingValueUnknown is a helper method that return true only in the case where the underlying value has a
// known type but the value is unknown. This method will return false if the underlying type is not known.
//
// IsUnknown should be used to determine if the dynamic value does not have a known type and the value is unknown.
//
// An example of a known type with an unknown underlying value would be:
//
//	types.DynamicValue(types.StringUnknown())
func (v DynamicValue) IsUnderlyingValueUnknown() bool {
	return v.value != nil && v.value.IsUnknown()
}
