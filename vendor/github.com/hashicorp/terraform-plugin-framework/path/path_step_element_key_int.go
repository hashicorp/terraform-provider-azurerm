// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "fmt"

// Ensure PathStepElementKeyInt satisfies the PathStep interface.
var _ PathStep = PathStepElementKeyInt(0)

// PathStepElementKeyInt is an attribute path transversal for an integer
// element of a list. List indexing starts a 0.
//
// Map elements must be transversed by PathStepElementKeyString.
// Object attributes must be transversed by PathStepAttributeName.
// Set elements must be transversed by PathStepElementKeyValue.
type PathStepElementKeyInt int64

// Equal returns true if the given PathStep is a PathStepAttributeName and the
// attribute name is equivalent.
func (s PathStepElementKeyInt) Equal(o PathStep) bool {
	other, ok := o.(PathStepElementKeyInt)

	if !ok {
		return false
	}

	return int64(s) == int64(other)
}

// ExpressionStep returns the ExpressionStep for the PathStep.
func (s PathStepElementKeyInt) ExpressionStep() ExpressionStep {
	return ExpressionStepElementKeyIntExact(s)
}

// String returns the human-readable representation of the element key.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s PathStepElementKeyInt) String() string {
	return fmt.Sprintf("[%d]", s)
}

// unexported satisfies the PathStep interface.
func (s PathStepElementKeyInt) unexported() {}
