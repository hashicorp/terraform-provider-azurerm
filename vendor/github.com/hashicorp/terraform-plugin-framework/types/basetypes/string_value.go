// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	_ StringValuable = StringValue{}
)

// StringValuable extends attr.Value for string value types.
// Implement this interface to create a custom String value type.
type StringValuable interface {
	attr.Value

	// ToStringValue should convert the value type to a String.
	ToStringValue(ctx context.Context) (StringValue, diag.Diagnostics)
}

// StringValuableWithSemanticEquals extends StringValuable with semantic
// equality logic.
type StringValuableWithSemanticEquals interface {
	StringValuable

	// StringSemanticEquals should return true if the given value is
	// semantically equal to the current value. This logic is used to prevent
	// Terraform data consistency errors and resource drift where a value change
	// may have inconsequential differences, such as spacing character removal
	// in JSON formatted strings.
	//
	// Only known values are compared with this method as changing a value's
	// state implicitly represents a different value.
	StringSemanticEquals(context.Context, StringValuable) (bool, diag.Diagnostics)
}

// NewStringNull creates a String with a null value. Determine whether the value is
// null via the String type IsNull method.
//
// Setting the deprecated String type Null, Unknown, or Value fields after
// creating a String with this function has no effect.
func NewStringNull() StringValue {
	return StringValue{
		state: attr.ValueStateNull,
	}
}

// NewStringUnknown creates a String with an unknown value. Determine whether the
// value is unknown via the String type IsUnknown method.
//
// Setting the deprecated String type Null, Unknown, or Value fields after
// creating a String with this function has no effect.
func NewStringUnknown() StringValue {
	return StringValue{
		state: attr.ValueStateUnknown,
	}
}

// NewStringValue creates a String with a known value. Access the value via the String
// type ValueString method.
//
// Setting the deprecated String type Null, Unknown, or Value fields after
// creating a String with this function has no effect.
func NewStringValue(value string) StringValue {
	return StringValue{
		state: attr.ValueStateKnown,
		value: value,
	}
}

// NewStringPointerValue creates a String with a null value if nil or a known
// value. Access the value via the String type ValueStringPointer method.
func NewStringPointerValue(value *string) StringValue {
	if value == nil {
		return NewStringNull()
	}

	return NewStringValue(*value)
}

// StringValue represents a UTF-8 string value.
type StringValue struct {
	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState

	// value contains the known value, if not null or unknown.
	value string
}

// Type returns a StringType.
func (s StringValue) Type(_ context.Context) attr.Type {
	return StringType{}
}

// ToTerraformValue returns the data contained in the *String as a tftypes.Value.
func (s StringValue) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	switch s.state {
	case attr.ValueStateKnown:
		if err := tftypes.ValidateValue(tftypes.String, s.value); err != nil {
			return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(tftypes.String, s.value), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(tftypes.String, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled String state in ToTerraformValue: %s", s.state))
	}
}

// Equal returns true if `other` is a String and has the same value as `s`.
func (s StringValue) Equal(other attr.Value) bool {
	o, ok := other.(StringValue)

	if !ok {
		return false
	}

	if s.state != o.state {
		return false
	}

	if s.state != attr.ValueStateKnown {
		return true
	}

	return s.value == o.value
}

// IsNull returns true if the String represents a null value.
func (s StringValue) IsNull() bool {
	return s.state == attr.ValueStateNull
}

// IsUnknown returns true if the String represents a currently unknown value.
func (s StringValue) IsUnknown() bool {
	return s.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the String value. Use
// the ValueString method for Terraform data handling instead.
//
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (s StringValue) String() string {
	if s.IsUnknown() {
		return attr.UnknownValueString
	}

	if s.IsNull() {
		return attr.NullValueString
	}

	return fmt.Sprintf("%q", s.value)
}

// ValueString returns the known string value. If String is null or unknown, returns
// "".
func (s StringValue) ValueString() string {
	return s.value
}

// ValueStringPointer returns a pointer to the known string value, nil for a
// null value, or a pointer to "" for an unknown value.
func (s StringValue) ValueStringPointer() *string {
	if s.IsNull() {
		return nil
	}

	return &s.value
}

// ToStringValue returns String.
func (s StringValue) ToStringValue(context.Context) (StringValue, diag.Diagnostics) {
	return s, nil
}
