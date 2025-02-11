// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure ExpressionStepElementKeyStringAny satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepElementKeyStringAny{}

// ExpressionStepElementKeyStringAny is an attribute path expression for a matching any
// string key within a map.
type ExpressionStepElementKeyStringAny struct{}

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyStringAny.
func (s ExpressionStepElementKeyStringAny) Equal(o ExpressionStep) bool {
	_, ok := o.(ExpressionStepElementKeyStringAny)

	return ok
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyStringAny condition.
func (s ExpressionStepElementKeyStringAny) Matches(pathStep PathStep) bool {
	_, ok := pathStep.(PathStepElementKeyString)

	return ok
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyStringAny) String() string {
	return `["*"]`
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyStringAny) unexported() {}
