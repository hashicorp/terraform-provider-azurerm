// Copyright (c) HashiCorp, Inc.
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

	// ErrPathInsideDynamicAttribute is used when AttributeAtPath is called on a path that doesn't
	// have a schema associated with it because it's nested in a dynamic attribute.
	ErrPathInsideDynamicAttribute = errors.New("path leads to element or attribute nested in a schema.DynamicAttribute")
)
