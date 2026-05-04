// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure ExpressionStepElementKeyIntAny satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepElementKeyIntAny{}

// ExpressionStepElementKeyIntAny is an attribute path expression for a matching any
// integer element key within a list.
type ExpressionStepElementKeyIntAny struct{}

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyIntAny.
func (s ExpressionStepElementKeyIntAny) Equal(o ExpressionStep) bool {
	_, ok := o.(ExpressionStepElementKeyIntAny)

	return ok
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyIntAny condition.
func (s ExpressionStepElementKeyIntAny) Matches(pathStep PathStep) bool {
	_, ok := pathStep.(PathStepElementKeyInt)

	return ok
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyIntAny) String() string {
	return "[*]"
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyIntAny) unexported() {}
