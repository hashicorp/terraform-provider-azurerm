package typehelpers

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/fwdiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.MapTypable        = (*mapObjectTypeOf[struct{}])(nil)
	_ NestedObjectCollectionType  = (*mapObjectTypeOf[struct{}])(nil)
	_ basetypes.MapValuable       = (*MapObjectValueOf[struct{}])(nil)
	_ NestedObjectCollectionValue = (*MapObjectValueOf[struct{}])(nil)
)

type mapObjectTypeOf[T any] struct {
	basetypes.MapType
}

func (m mapObjectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objectTypeNewObjectPtr[T](ctx)
}

func (m mapObjectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	return NewMapObjectValueOfNull[T](ctx), diags
}

func (m mapObjectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, a any) (attr.Value, diag.Diagnostics) {
	// TODO implement me
	panic("implement me")
}

func (m mapObjectTypeOf[T]) NewObjectSlice(ctx context.Context, i int, i2 int) (any, diag.Diagnostics) {
	// TODO implement me
	panic("implement me")
}

func (m mapObjectTypeOf[T]) ValueFromObjectSlice(ctx context.Context, a any) (attr.Value, diag.Diagnostics) {
	// TODO implement me
	panic("implement me")
}

type MapObjectValueOf[T any] struct {
	basetypes.MapValue
}

func (m MapObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	// TODO implement me
	panic("implement me")
}

func (m MapObjectValueOf[T]) ToObjectSlice(ctx context.Context) (any, diag.Diagnostics) {
	// TODO implement me
	panic("implement me")
}

func NewMapValueObjectOfUnknown[T any](ctx context.Context) MapObjectValueOf[T] {
	return MapObjectValueOf[T]{MapValue: basetypes.NewMapUnknown(NewObjectTypeOf[T](ctx))}
}

func NewMapObjectValueOf[T any](ctx context.Context, t map[string]T) (MapObjectValueOf[T], diag.Diagnostics) {
	return newMapObjectValueOf[T](ctx, t)
}

func NewMapObjectValueOfMust[T any](ctx context.Context, t map[string]T) MapObjectValueOf[T] {
	return fwdiag.Must(newMapObjectValueOf[T](ctx, t))
}

func newMapObjectValueOf[T any](ctx context.Context, elements any) (MapObjectValueOf[T], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	typ, d := newObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapValueObjectOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewMapValueFrom(ctx, typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapValueObjectOfUnknown[T](ctx), diags
	}

	return MapObjectValueOf[T]{MapValue: v}, diags
}

func NewMapObjectValueOfNull[T any](ctx context.Context) MapObjectValueOf[T] {
	return MapObjectValueOf[T]{MapValue: basetypes.NewMapNull(NewObjectTypeOf[T](ctx))}
}
