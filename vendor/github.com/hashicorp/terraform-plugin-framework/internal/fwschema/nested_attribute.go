// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

// NestedAttribute defines a schema attribute that contains nested attributes.
type NestedAttribute interface {
	Attribute

	// GetNestedObject should return the object underneath the nested
	// attribute. For single nesting mode, the NestedAttributeObject can be
	// generated from the Attribute.
	GetNestedObject() NestedAttributeObject

	// GetNestingMode should return the nesting mode (list, map, set, or
	// single) of the nested attributes or left unset if this Attribute
	// does not represent nested attributes.
	GetNestingMode() NestingMode
}

// NestedAttributesEqual is a helper function to perform equality testing on two
// NestedAttribute. NestedAttribute Equal implementations should still compare
// the concrete types in addition to using this helper.
func NestedAttributesEqual(a, b NestedAttribute) bool {
	if !AttributesEqual(a, b) {
		return false
	}

	if a.GetNestingMode() != b.GetNestingMode() {
		return false
	}

	return a.GetNestedObject().Equal(b.GetNestedObject())
}
