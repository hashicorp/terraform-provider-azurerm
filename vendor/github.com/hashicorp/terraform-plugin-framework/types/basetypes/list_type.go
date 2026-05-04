// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ ListTypable = ListType{}

// ListTypable extends attr.Type for list types.
// Implement this interface to create a custom ListType type.
type ListTypable interface {
	attr.Type

	// ValueFromList should convert the List to a ListValuable type.
	ValueFromList(context.Context, ListValue) (ListValuable, diag.Diagnostics)
}

// ListType is an AttributeType representing a list of values. All values must
// be of the same type, which the provider must specify as the ElemType
// property.
type ListType struct {
	ElemType attr.Type
}

// ElementType returns the attr.Type elements will be created from.
func (l ListType) ElementType() attr.Type {
	if l.ElemType == nil {
		return missingType{}
	}

	return l.ElemType
}

// WithElementType returns a ListType that is identical to `l`, but with the
// element type set to `typ`.
func (l ListType) WithElementType(typ attr.Type) attr.TypeWithElementType {
	return ListType{ElemType: typ}
}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the AttributeType to something Terraform
// can understand.
func (l ListType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.List{
		ElementType: l.ElementType().TerraformType(ctx),
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value.
// This is meant to convert the tftypes.Value into a more convenient Go
// type for the provider to consume the data with.
func (l ListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewListNull(l.ElementType()), nil
	}

	// MAINTAINER NOTE:
	// ListType does not support DynamicType as an element type. It is not explicitly prevented from being created with the
	// Framework type system, but the Framework-supported ListAttribute, ListNestedAttribute, and ListNestedBlock all prevent DynamicType
	// from being used as an element type. An attempt to use DynamicType as the element type will eventually lead you to an error on this line :)
	//
	// In the future, if we ever need to support a list of dynamic element types, this type equality check will need to be modified to allow
	// dynamic types to not return an error, as the tftypes.Value coming in (if known) will be a concrete value, for example:
	//
	// - l.TerraformType(ctx): tftypes.List[tftypes.DynamicPseudoType]
	// - in.Type(): tftypes.List[tftypes.String]
	//
	// The `ValueFromTerraform` function for a dynamic type will be able create the correct concrete dynamic value with this modification in place.
	//
	if !in.Type().Equal(l.TerraformType(ctx)) {
		return nil, fmt.Errorf("can't use %s as value of List with ElementType %T, can only use %s values", in.String(), l.ElementType(), l.ElementType().TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		return NewListUnknown(l.ElementType()), nil
	}
	if in.IsNull() {
		return NewListNull(l.ElementType()), nil
	}
	val := []tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make([]attr.Value, 0, len(val))
	for _, elem := range val {
		av, err := l.ElementType().ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, av)
	}
	// ValueFromTerraform above on each element should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewListValueMust(l.ElementType(), elems), nil
}

// Equal returns true if `o` is also a ListType and has the same ElemType.
func (l ListType) Equal(o attr.Type) bool {
	// Preserve prior ElemType nil check behavior
	if l.ElementType().Equal(missingType{}) {
		return false
	}

	other, ok := o.(ListType)

	if !ok {
		return false
	}

	return l.ElementType().Equal(other.ElementType())
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// list.
func (l ListType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if _, ok := step.(tftypes.ElementKeyInt); !ok {
		return nil, fmt.Errorf("cannot apply step %T to ListType", step)
	}

	return l.ElementType(), nil
}

// String returns a human-friendly description of the ListType.
func (l ListType) String() string {
	return "types.ListType[" + l.ElementType().String() + "]"
}

// Validate validates all elements of the list that are of type
// xattr.TypeWithValidate.
func (l ListType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.List{}) {
		err := fmt.Errorf("expected List value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"List Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var elems []tftypes.Value

	if err := in.As(&elems); err != nil {
		diags.AddAttributeError(
			path,
			"List Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
	validatableType, isValidatable := l.ElementType().(xattr.TypeWithValidate)
	if !isValidatable {
		return diags
	}

	for index, elem := range elems {
		if !elem.IsFullyKnown() {
			continue
		}
		diags = append(diags, validatableType.Validate(ctx, elem, path.AtListIndex(index))...)
	}

	return diags
}

// ValueType returns the Value type.
func (l ListType) ValueType(_ context.Context) attr.Value {
	return ListValue{
		elementType: l.ElementType(),
	}
}

// ValueFromList returns a ListValuable type given a List.
func (l ListType) ValueFromList(_ context.Context, list ListValue) (ListValuable, diag.Diagnostics) {
	return list, nil
}
