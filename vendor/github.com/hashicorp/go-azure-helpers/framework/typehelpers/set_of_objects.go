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

type setNestedObjectTypeOf[T any] struct {
	basetypes.SetType
}

type SetNestedObjectValueOf[T any] struct {
	basetypes.SetValue
}

var (
	_ basetypes.SetTypable        = (*setNestedObjectTypeOf[struct{}])(nil)
	_ NestedObjectCollectionType  = (*setNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.SetValuable       = (*SetNestedObjectValueOf[struct{}])(nil)
	_ NestedObjectCollectionValue = (*SetNestedObjectValueOf[struct{}])(nil)
)

// NewSetNestedObjectTypeOf returns a wrapped `ListType` of a given type `T`
// This can be used in provider resource implementations for Custom Types (i.e. sets of objects, which is to say
// "blocks")
func NewSetNestedObjectTypeOf[T any](ctx context.Context) setNestedObjectTypeOf[T] {
	return setNestedObjectTypeOf[T]{basetypes.SetType{ElemType: NewObjectTypeOf[T](ctx)}}
}

func (s setNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(setNestedObjectTypeOf[T])

	if !ok {
		return false
	}

	return s.SetType.Equal(other.SetType)
}

func (s setNestedObjectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objectTypeNewObjectPtr[T](ctx)
}

func (s setNestedObjectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	return NewSetNestedObjectValueOfNull[T](ctx), diags
}

func (s setNestedObjectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, ptr any) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if v, ok := ptr.(*T); ok {
		v, d := NewSetNestedObjectValueOfPtr(ctx, v)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: want %T, got %T", (*T)(nil), ptr)))
	return nil, diags
}

func (s setNestedObjectTypeOf[T]) ValueFromObjectSlice(ctx context.Context, slice any) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if v, ok := slice.([]*T); ok {
		return NewSetNestedObjectValueOfSlice(ctx, v)
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid slice value", fmt.Sprintf("incorrect type: want %T, got %T", (*[]T)(nil), slice)))
	return nil, diags
}

func (s setNestedObjectTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewSetNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	typ, d := newObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewSetValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{SetValue: v}, diags
}

func (s setNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := s.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := s.ValueFromSet(ctx, setValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

func (s setNestedObjectTypeOf[T]) ValueType(_ context.Context) attr.Value {
	return SetNestedObjectValueOf[T]{}
}

func (s setNestedObjectTypeOf[T]) NewObjectSlice(ctx context.Context, l int, cap int) (any, diag.Diagnostics) {
	return nestedObjectTypeNewObjectSlice[T](ctx, l, cap)
}

func (s SetNestedObjectValueOf[T]) ToObjectSlice(ctx context.Context) (any, diag.Diagnostics) {
	return s.ToSlice(ctx)
}

func (s SetNestedObjectValueOf[T]) ToSlice(ctx context.Context) ([]*T, diag.Diagnostics) {
	return nestedObjectValueObjectSlice[T](ctx, s.SetValue)
}

func (s SetNestedObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return s.ToPtr(ctx)
}

func (s SetNestedObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return nestedObjectValueObjectPtr[T](ctx, s.SetValue)
}

func NewSetNestedObjectValueOfNull[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	return SetNestedObjectValueOf[T]{SetValue: basetypes.NewSetNull(NewObjectTypeOf[T](ctx))}
}

func NewSetNestedObjectValueOfPtr[T any](ctx context.Context, t *T) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	return NewSetNestedObjectValueOfSlice(ctx, []*T{t})
}

func NewSetNestedObjectValueOfSlice[T any](ctx context.Context, ts []*T) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	return newSetNestedObjectValueOf[T](ctx, ts)
}

func NewSetNestedObjectValueOfValueSlice[T any](ctx context.Context, ts []T) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	return newSetNestedObjectValueOf[T](ctx, ts)
}

func NewSetNestedObjectValueOfValueSliceMust[T any](ctx context.Context, ts []T) SetNestedObjectValueOf[T] {
	return fwdiag.Must(NewSetNestedObjectValueOfValueSlice[T](ctx, ts))
}

func newSetNestedObjectValueOf[T any](ctx context.Context, elements any) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	typ, d := newObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewSetValueFrom(ctx, typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewSetNestedObjectValueOfUnknown[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{SetValue: v}, diags
}

func NewSetNestedObjectValueOfUnknown[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	return SetNestedObjectValueOf[T]{SetValue: basetypes.NewSetUnknown(NewObjectTypeOf[T](ctx))}
}

func (s SetNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetNestedObjectValueOf[T])

	if !ok {
		return false
	}

	return s.SetValue.Equal(other.SetValue)
}

func (s SetNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewSetNestedObjectTypeOf[T](ctx)
}

func (s setNestedObjectTypeOf[T]) String() string {
	var t T
	return fmt.Sprintf("SetNestedObjectTypeOf[%T]", t)
}
