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
	_ basetypes.SetTypable  = (*setTypeOf[basetypes.StringValue])(nil)
	_ basetypes.SetValuable = (*SetValueOf[basetypes.StringValue])(nil)
)

type setTypeOf[T attr.Value] struct {
	basetypes.SetType
}

var (
	SetOfStringType  = setTypeOf[basetypes.StringValue]{basetypes.SetType{ElemType: basetypes.StringType{}}}
	SetOfInt64Type   = setTypeOf[basetypes.Int64Value]{basetypes.SetType{ElemType: basetypes.Int64Type{}}}
	SetOfFloat64Type = setTypeOf[basetypes.Float64Value]{basetypes.SetType{ElemType: basetypes.Float64Type{}}}
)

func NewSetTypeOf[T attr.Value](ctx context.Context) setTypeOf[T] {
	return newSetTypeOf[T](ctx)
}

func newSetTypeOf[T attr.Value](ctx context.Context) setTypeOf[T] {
	return setTypeOf[T]{basetypes.SetType{ElemType: newAttrTypeOf[T](ctx)}}
}

func (t setTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(setTypeOf[T])

	if !ok {
		return false
	}

	return t.SetType.Equal(other.SetType)
}

func (t setTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("SetTypeOf[%T]", zero)
}

func (t setTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewSetValueOfNull[T](ctx), diags
	}

	if in.IsUnknown() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewSetValue(newAttrTypeOf[T](ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: v}, diags
}

func (t setTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

func (t setTypeOf[T]) ValueType(_ context.Context) attr.Value {
	return SetValueOf[T]{}
}

type SetValueOf[T attr.Value] struct {
	basetypes.SetValue
}

func (v SetValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetValueOf[T])

	if !ok {
		return false
	}

	return v.SetValue.Equal(other.SetValue)
}

func (v SetValueOf[T]) Type(ctx context.Context) attr.Type {
	return newSetTypeOf[T](ctx)
}

func NewSetValueOfNull[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetNull(newAttrTypeOf[T](ctx))}
}

func NewSetValueOfUnknown[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetUnknown(newAttrTypeOf[T](ctx))}
}

func NewSetValueOf[T attr.Value](ctx context.Context, elements []attr.Value) (SetValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewSetValue(newAttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetValueOfUnknown[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: v}, diags
}

func NewSetValueOfMust[T attr.Value](ctx context.Context, elements []attr.Value) SetValueOf[T] {
	return fwdiag.Must(NewSetValueOf[T](ctx, elements))
}
