// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure ExpressionStepAttributeNameExact satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepAttributeNameExact("")

// ExpressionStepAttributeNameExact is an attribute path expression for an
// exact attribute name match within an object.
type ExpressionStepAttributeNameExact string

// Equal returns true if the given ExpressionStep is a
// ExpressionStepAttributeNameExact and the attribute name is equivalent.
func (s ExpressionStepAttributeNameExact) Equal(o ExpressionStep) bool {
	other, ok := o.(ExpressionStepAttributeNameExact)

	if !ok {
		return false
	}

	return string(s) == string(other)
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepAttributeNameExact condition.
func (s ExpressionStepAttributeNameExact) Matches(pathStep PathStep) bool {
	pathStepAttributeName, ok := pathStep.(PathStepAttributeName)

	if !ok {
		return false
	}

	return string(s) == string(pathStepAttributeName)
}

// String returns the human-readable representation of the attribute name
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepAttributeNameExact) String() string {
	return string(s)
}

// unexported satisfies the Step interface.
func (s ExpressionStepAttributeNameExact) unexported() {}
