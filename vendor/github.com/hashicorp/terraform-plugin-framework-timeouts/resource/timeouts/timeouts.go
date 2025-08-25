// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ basetypes.ObjectTypable  = Type{}
	_ basetypes.ObjectValuable = Value{}
)

// Type is an attribute type that represents timeouts.
type Type struct {
	basetypes.ObjectType
}

// String returns a human-readable representation of the type.
func (t Type) String() string {
	return "timeouts.Type"
}

// ValueFromObject returns a Value given a basetypes.ObjectValue.
func (t Type) ValueFromObject(_ context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	value := Value{
		Object: in,
	}

	return value, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
// Value embeds the types.Object value returned from calling ValueFromTerraform on the
// types.ObjectType embedded in Type.
func (t Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := t.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	obj, ok := val.(types.Object)
	if !ok {
		return nil, fmt.Errorf("%T cannot be used as types.Object", val)
	}

	return Value{
		obj,
	}, err
}

// ValueType returns the associated Value type for debugging.
func (t Type) ValueType(context.Context) attr.Value {
	// It does not need to be a fully valid implementation of the type.
	return Value{}
}

// Equal returns true if `candidate` is also a Type and has the same
// AttributeTypes.
func (t Type) Equal(candidate attr.Type) bool {
	other, ok := candidate.(Type)
	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

// Value represents an object containing values to be used as time.Duration for timeouts.
type Value struct {
	types.Object
}

// Equal returns true if the Value is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (t Value) Equal(c attr.Value) bool {
	other, ok := c.(Value)

	if !ok {
		return false
	}

	return t.Object.Equal(other.Object)
}

// ToObjectValue returns the underlying ObjectValue.
func (v Value) ToObjectValue(_ context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return v.Object, nil
}

// Type returns a Type with the same attribute types as `t`.
func (t Value) Type(ctx context.Context) attr.Type {
	return Type{
		types.ObjectType{
			AttrTypes: t.AttributeTypes(ctx),
		},
	}
}

// Create attempts to retrieve the "create" attribute and parse it as time.Duration.
// If any diagnostics are generated they are returned along with the supplied default timeout.
func (t Value) Create(ctx context.Context, defaultTimeout time.Duration) (time.Duration, diag.Diagnostics) {
	return t.getTimeout(ctx, attributeNameCreate, defaultTimeout)
}

// Read attempts to retrieve the "read" attribute and parse it as time.Duration.
// If any diagnostics are generated they are returned along with the supplied default timeout.
func (t Value) Read(ctx context.Context, defaultTimeout time.Duration) (time.Duration, diag.Diagnostics) {
	return t.getTimeout(ctx, attributeNameRead, defaultTimeout)
}

// Update attempts to retrieve the "update" attribute and parse it as time.Duration.
// If any diagnostics are generated they are returned along with the supplied default timeout.
func (t Value) Update(ctx context.Context, defaultTimeout time.Duration) (time.Duration, diag.Diagnostics) {
	return t.getTimeout(ctx, attributeNameUpdate, defaultTimeout)
}

// Delete attempts to retrieve the "delete" attribute and parse it as time.Duration.
// If any diagnostics are generated they are returned along with the supplied default timeout.
func (t Value) Delete(ctx context.Context, defaultTimeout time.Duration) (time.Duration, diag.Diagnostics) {
	return t.getTimeout(ctx, attributeNameDelete, defaultTimeout)
}

func (t Value) getTimeout(ctx context.Context, timeoutName string, defaultTimeout time.Duration) (time.Duration, diag.Diagnostics) {
	var diags diag.Diagnostics

	value, ok := t.Object.Attributes()[timeoutName]
	if !ok {
		tflog.Info(ctx, timeoutName+" timeout configuration not found, using provided default")

		return defaultTimeout, diags
	}

	if value.IsNull() || value.IsUnknown() {
		tflog.Info(ctx, timeoutName+" timeout configuration is null or unknown, using provided default")

		return defaultTimeout, diags
	}

	// No type assertion check is required as the schema guarantees that the object attributes
	// are types.String.
	timeout, err := time.ParseDuration(value.(types.String).ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic(
			"Timeout Cannot Be Parsed",
			fmt.Sprintf("timeout for %q cannot be parsed, %s", timeoutName, err),
		))

		return defaultTimeout, diags
	}

	return timeout, diags
}
