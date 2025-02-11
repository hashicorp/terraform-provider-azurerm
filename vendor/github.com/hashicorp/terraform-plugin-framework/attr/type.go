// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Type defines an interface for describing a kind of attribute. Types are
// collections of constraints and behaviors such that they can be reused on
// multiple attributes easily.
//
// Refer also to the xattr package, which contains additional extensions for
// Type, such as validation.
type Type interface {
	// TerraformType returns the tftypes.Type that should be used to
	// represent this type. This constrains what user input will be
	// accepted and what kind of data can be set in state. The framework
	// will use this to translate the Type to something Terraform can
	// understand.
	TerraformType(context.Context) tftypes.Type

	// ValueFromTerraform returns a Value given a tftypes.Value. This is
	// meant to convert the tftypes.Value into a more convenient Go type
	// for the provider to consume the data with.
	ValueFromTerraform(context.Context, tftypes.Value) (Value, error)

	// ValueType should return the attr.Value type returned by
	// ValueFromTerraform. The returned attr.Value can be any null, unknown,
	// or known value for the type, as this is intended for type detection
	// and improving error diagnostics.
	ValueType(context.Context) Value

	// Equal should return true if the Type is considered equivalent to the
	// Type passed as an argument.
	//
	// Most types should verify the associated Type is exactly equal to prevent
	// potential data consistency issues. For example:
	//
	//  - basetypes.Number is inequal to basetypes.Int64 or basetypes.Float64
	//  - basetypes.String is inequal to a custom Go type that embeds it
	//
	Equal(Type) bool

	// String should return a human-friendly version of the Type.
	String() string

	tftypes.AttributePathStepper
}

// TypeWithAttributeTypes extends the Type interface to include information about
// attribute types. Attribute types are part of the definition of an object type.
type TypeWithAttributeTypes interface {
	Type

	// WithAttributeTypes returns a new copy of the type with its
	// attribute types set.
	WithAttributeTypes(map[string]Type) TypeWithAttributeTypes

	// AttributeTypes returns the object's attribute types.
	AttributeTypes() map[string]Type
}

// TypeWithElementType extends the Type interface to include information about the type
// all elements will share. Element types are part of the definition of a list,
// set, or map type.
type TypeWithElementType interface {
	Type

	// WithElementType returns a new copy of the type with its element type
	// set.
	WithElementType(Type) TypeWithElementType

	// ElementType returns the type's element type.
	ElementType() Type
}

// TypeWithElementTypes extends the Type interface to include information about the
// types of each element. Element types are part of the definition of a tuple
// type.
type TypeWithElementTypes interface {
	Type

	// WithElementTypes returns a new copy of the type with its elements'
	// types set.
	WithElementTypes([]Type) TypeWithElementTypes

	// ElementTypes returns the type's elements' types.
	ElementTypes() []Type
}

// TypeWithPlaintextDescription extends the Type interface to include a
// Description method, used to bundle extra information to include in attribute
// descriptions with the Type. It expects the description to be written as
// plain text, with no special formatting.
type TypeWithPlaintextDescription interface {
	Type

	// Description returns a practitioner-friendly explanation of the type
	// and the constraints of the data it accepts and returns. It will be
	// combined with the Description associated with the Attribute.
	Description(context.Context) string
}

// TypeWithMarkdownDescription extends the Type interface to include a
// MarkdownDescription method, used to bundle extra information to include in
// attribute descriptions with the Type. It expects the description to be
// formatted for display with Markdown.
type TypeWithMarkdownDescription interface {
	Type

	// MarkdownDescription returns a practitioner-friendly explanation of
	// the type and the constraints of the data it accepts and returns. It
	// will be combined with the MarkdownDescription associated with the
	// Attribute.
	MarkdownDescription(context.Context) string
}
