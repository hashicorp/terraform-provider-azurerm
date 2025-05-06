package typehelpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type NestedObjectMapValue interface {
	NestedObjectValue

	// ToObjectMap returns the Value as an object pointer (
	ToObjectMap(context.Context) (any, diag.Diagnostics)
}

type MapObjectType interface {
	attr.Type

	// NewObjectPtr returns a new, empty value as an object pointer (Go *struct).
	NewObjectPtr(context.Context) (any, diag.Diagnostics)

	// NullValue returns a Null Value.
	NullValue(context.Context) (attr.Value, diag.Diagnostics)

	// ValueFromObjectPtr returns a Value given an object pointer (Go *struct).
	ValueFromObjectPtr(context.Context, any) (attr.Value, diag.Diagnostics)
}

type MapObjectCollectionType interface {
	MapObjectType

	// NewObjectMap returns a new value as an object map (Go map[string]struct).
	NewObjectMap(context.Context, int, int) (any, diag.Diagnostics)

	// ValueFromObjectMap returns a Value given an object pointer (Go map[string]struct).
	ValueFromObjectMap(context.Context, any) (attr.Value, diag.Diagnostics)
}

type MapObjectValue interface {
	attr.Value

	// ToObjectMapPtr returns the value as an object pointer (Go *struct).
	ToObjectMapPtr(context.Context) (any, diag.Diagnostics)
}

// MapObjectCollectionValue extends the NestedObjectValue interface for values that represent collections of nested Objects.
type MapObjectCollectionValue interface {
	NestedObjectValue

	// ToObjectMap returns the value as an object slice (Go []*struct).
	ToObjectMap(context.Context) (map[string]any, diag.Diagnostics)
}

// // mapValueWithElements extends the Value interface for values that have an Elements method.
// type mapValueWithElements interface {
// 	attr.Value
//
// 	Elements() map[string]attr.Value
// }
