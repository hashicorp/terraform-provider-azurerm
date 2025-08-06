// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure ExpressionStepElementKeyValueAny satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepElementKeyValueAny{}

// ExpressionStepElementKeyValueAny is an attribute path expression for a matching any
// Value element within a set.
type ExpressionStepElementKeyValueAny struct{}

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyValueAny.
func (s ExpressionStepElementKeyValueAny) Equal(o ExpressionStep) bool {
	_, ok := o.(ExpressionStepElementKeyValueAny)

	return ok
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyValueAny condition.
func (s ExpressionStepElementKeyValueAny) Matches(pathStep PathStep) bool {
	_, ok := pathStep.(PathStepElementKeyValue)

	return ok
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyValueAny) String() string {
	return "[Value(*)]"
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyValueAny) unexported() {}
