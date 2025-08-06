// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure PathStepAttributeName satisfies the PathStep interface.
var _ PathStep = PathStepAttributeName("")

// PathStepAttributeName is an attribute path tranversal for an attribute name
// within an object.
//
// List elements must be transversed by PathStepElementKeyInt.
// Map elements must be transversed by PathStepElementKeyString.
// Set elements must be transversed by PathStepElementKeyValue.
type PathStepAttributeName string

// Equal returns true if the given PathStep is a PathStepAttributeName and the
// attribute name is equivalent.
func (s PathStepAttributeName) Equal(o PathStep) bool {
	other, ok := o.(PathStepAttributeName)

	if !ok {
		return false
	}

	return string(s) == string(other)
}

// ExpressionStep returns the ExpressionStep for the PathStep.
func (s PathStepAttributeName) ExpressionStep() ExpressionStep {
	return ExpressionStepAttributeNameExact(s)
}

// String returns the human-readable representation of the attribute name.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s PathStepAttributeName) String() string {
	return string(s)
}

// unexported satisfies the PathStep interface.
func (s PathStepAttributeName) unexported() {}
