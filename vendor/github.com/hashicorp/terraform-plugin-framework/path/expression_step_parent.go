// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// Ensure StepParent satisfies the ExpressionStep interface.
var _ ExpressionStep = ExpressionStepParent{}

// StepParent is an attribute path expression for a traversing to the parent
// attribute path relative to the current one. This is intended only for the
// start of attribute-level expressions which will be combined with the current
// attribute path being called.
type ExpressionStepParent struct{}

// Equal returns true if the given ExpressionStep is a ExpressionStepParent.
func (s ExpressionStepParent) Equal(o ExpressionStep) bool {
	_, ok := o.(ExpressionStepParent)

	return ok
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepParent condition.
func (s ExpressionStepParent) Matches(_ PathStep) bool {
	// This return value should have no effect, as this Step is a
	// sentinel type, rather than one that should be used in matching.
	return false
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepParent) String() string {
	return "<"
}

// unexported satisfies the Step interface.
func (s ExpressionStepParent) unexported() {}
