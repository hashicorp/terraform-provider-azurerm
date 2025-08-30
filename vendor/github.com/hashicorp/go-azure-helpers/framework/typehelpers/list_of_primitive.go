// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/fwdiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ListTypable  = (*listTypeOf[basetypes.StringValue])(nil)
	_ basetypes.ListValuable = (*ListValueOf[basetypes.StringValue])(nil)
)

type listTypeOf[T attr.Value] struct {
	basetypes.ListType
}

var (
	ListOfStringType  = listTypeOf[basetypes.StringValue]{basetypes.ListType{ElemType: basetypes.StringType{}}}
	ListOfInt64Type   = listTypeOf[basetypes.Int64Value]{basetypes.ListType{ElemType: basetypes.Int64Type{}}}
	ListOfFloat64Type = listTypeOf[basetypes.Float64Value]{basetypes.ListType{ElemType: basetypes.Float64Type{}}}
)

func NewListTypeOf[T attr.Value](ctx context.Context) listTypeOf[T] {
	return newListTypeOf[T](ctx)
}

func newListTypeOf[T attr.Value](ctx context.Context) listTypeOf[T] {
	return listTypeOf[T]{basetypes.ListType{ElemType: newAttrTypeOf[T](ctx)}}
}

func (t listTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(listTypeOf[T])

	if !ok {
		return false
	}

	return t.ListType.Equal(other.ListType)
}

func (t listTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("ListTypeOf[%T]", zero)
}

func (t listTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewListValueOfNull[T](ctx), diags
	}

	if in.IsUnknown() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewListValue(newAttrTypeOf[T](ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
}

func (t listTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func (t listTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListValueOf[T]{}
}

type ListValueOf[T attr.Value] struct {
	basetypes.ListValue
}

func (v ListValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListValueOf[T])

	if !ok {
		return false
	}

	return v.ListValue.Equal(other.ListValue)
}

func (v ListValueOf[T]) Type(ctx context.Context) attr.Type {
	return newListTypeOf[T](ctx)
}

func NewListValueOfNull[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListNull(newAttrTypeOf[T](ctx))}
}

func NewListValueOfUnknown[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListUnknown(newAttrTypeOf[T](ctx))}
}

func NewListValueOf[T attr.Value](ctx context.Context, elements []attr.Value) (ListValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewListValue(newAttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
}

func NewListValueOfMust[T attr.Value](ctx context.Context, elements []attr.Value) ListValueOf[T] {
	return fwdiag.Must(NewListValueOf[T](ctx, elements))
}
