// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-azure-helpers/framework/fwdiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ObjectTypable  = (*objectTypeOf[struct{}])(nil)
	_ NestedObjectType         = (*objectTypeOf[struct{}])(nil)
	_ basetypes.ObjectValuable = (*ObjectValueOf[struct{}])(nil)
	_ NestedObjectValue        = (*ObjectValueOf[struct{}])(nil)
)

type objectTypeOf[T any] struct {
	basetypes.ObjectType
}

func newObjectTypeOf[T any](ctx context.Context) (objectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return objectTypeOf[T]{}, diags
	}

	return objectTypeOf[T]{basetypes.ObjectType{AttrTypes: m}}, diags
}

func NewObjectTypeOf[T any](ctx context.Context) objectTypeOf[T] {
	return fwdiag.Must(newObjectTypeOf[T](ctx))
}

func (o objectTypeOf[T]) Equal(t attr.Type) bool {
	other, ok := t.(objectTypeOf[T])
	if !ok {
		return false
	}

	return o.ObjectType.Equal(other.ObjectType)
}

func (o objectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("ObjectTypeOf[%T]", zero)
}

func (o objectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := o.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	objectValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	objectValuable, diags := o.ValueFromObject(ctx, objectValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ObjectValue to ObjectValuable: %v", diags)
	}

	return objectValuable, nil
}

func (o objectTypeOf[T]) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	m, d := AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewObjectValue(m, in.Attributes())
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	value := ObjectValueOf[T]{
		ObjectValue: v,
	}

	return value, diags
}

func (objectTypeOf[T]) ValueType(_ context.Context) attr.Value {
	return ObjectValueOf[T]{}
}

func (o objectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objectTypeNewObjectPtr[T](ctx)
}

func (o objectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	return NewObjectValueOfNull[T](ctx), diags
}

func (o objectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, a any) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if v, ok := a.(*T); ok {
		v, d := NewObjectValueOf(ctx, v)
		diags.Append(d...)
		return v, diags
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: expected %T, got %T", (*T)(nil), a)))
	return nil, diags
}

func objectTypeNewObjectPtr[T any](ctx context.Context) (*T, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	t := new(T)

	diags.Append(NullOutObjectPtrFields(ctx, t)...)
	if diags.HasError() {
		return nil, diags
	}

	return t, diags
}

func objectValueObjectPtr[T any](ctx context.Context, val attr.Value) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	ptr, d := objectTypeNewObjectPtr[T](ctx)
	debugType := reflect.TypeOf(val)
	d.Append(diag.NewWarningDiagnostic(debugType.String(), "ObjectValue"))
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	diags.Append(val.(ObjectValueOf[T]).As(ctx, ptr, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	return ptr, diags
}

func NullOutObjectPtrFields[T any](ctx context.Context, t *T) diag.Diagnostics {
	var diags diag.Diagnostics
	val := reflect.ValueOf(t)
	typ := val.Type().Elem()

	if typ.Kind() != reflect.Struct {
		return diags
	}

	val = val.Elem()

	for i := 0; i < typ.NumField(); i++ {
		val := val.Field(i)
		if !val.CanInterface() {
			continue
		}

		attrValue, err := NullValueOf(ctx, val.Interface())
		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("attr.Type.ValueFromTerraform", err.Error()))
			return diags
		}

		if attrValue == nil {
			continue
		}

		val.Set(reflect.ValueOf(attrValue))
	}

	return diags
}

// ObjectValueOf represents a Terraform Plugin Framework Object value whose corresponding Go type is the structure T.
type ObjectValueOf[T any] struct {
	basetypes.ObjectValue
}

func (v ObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ObjectValueOf[T])
	if !ok {
		return false
	}

	return v.ObjectValue.Equal(other.ObjectValue)
}

func (v ObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewObjectTypeOf[T](ctx)
}

func (v ObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToPtr(ctx)
}

func (v ObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return objectValueObjectPtr[T](ctx, v)
}

func NewObjectValueOfNull[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectNull(AttributeTypesMust[T](ctx))}
}

func NewObjectValueOfUnknown[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectUnknown(AttributeTypesMust[T](ctx))}
}

func NewObjectValueOf[T any](ctx context.Context, t *T) (ObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewObjectValueFrom(ctx, m, t)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	return ObjectValueOf[T]{ObjectValue: v}, diags
}

func nestedObjectTypeNewObjectSlice[T any](_ context.Context, l, cap int) ([]*T, diag.Diagnostics) { //nolint:unparam
	var diags diag.Diagnostics

	return make([]*T, l, cap), diags
}
