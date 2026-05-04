// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Ensure PathStepElementKeyValue satisfies the PathStep interface.
// var _ PathStep = PathStepElementKeyValue(/* ... */)

// PathStepElementKeyValue is an attribute path transversal for a Value element
// of a set. Sets do not use integer-based indexing.
//
// List elements must be transversed by PathStepElementKeyInt.
// Map elements must be transversed by PathStepElementKeyString.
// Object attributes must be transversed by PathStepAttributeName.
type PathStepElementKeyValue struct {
	// Value is an interface, so it cannot be type aliased with methods.
	attr.Value
}

// Equal returns true if the given PathStep is a PathStepAttributeName and the
// attribute name is equivalent.
func (s PathStepElementKeyValue) Equal(o PathStep) bool {
	other, ok := o.(PathStepElementKeyValue)

	if !ok {
		return false
	}

	return s.Value.Equal(other.Value)
}

// ExpressionStep returns the ExpressionStep for the PathStep.
func (s PathStepElementKeyValue) ExpressionStep() ExpressionStep {
	return ExpressionStepElementKeyValueExact(s)
}

// String returns the human-readable representation of the element key.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s PathStepElementKeyValue) String() string {
	return fmt.Sprintf("[Value(%s)]", s.Value.String())
}

// unexported satisfies the PathStep interface.
func (s PathStepElementKeyValue) unexported() {}
