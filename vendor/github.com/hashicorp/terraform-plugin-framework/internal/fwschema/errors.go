// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwschema

import "errors"

var (
	// ErrPathInsideAtomicAttribute is used when AttributeAtPath is called
	// on a path that doesn't have a schema associated with it, because
	// it's an element, attribute, or block of a complex type, not a nested
	// attribute.
	ErrPathInsideAtomicAttribute = errors.New("path leads to element, attribute, or block of a schema.Attribute that has no schema associated with it")

	// ErrPathIsBlock is used when AttributeAtPath is called on a path is a
	// block, not an attribute. Use blockAtPath on the path instead.
	ErrPathIsBlock = errors.New("path leads to block, not an attribute")

	// ErrPathIsAttribute is used when BlockAtPath is called on a path that is
	// an attribute, not a block. Use attributeAtPath on the path instead.
	ErrPathIsAttribute = errors.New("path leads to attribute, not a block")

	// ErrPathInsideDynamicAttribute is used when AttributeAtPath is called on a path that doesn't
	// have a schema associated with it because it's nested in a dynamic attribute.
	ErrPathInsideDynamicAttribute = errors.New("path leads to element or attribute nested in a schema.DynamicAttribute")
)
