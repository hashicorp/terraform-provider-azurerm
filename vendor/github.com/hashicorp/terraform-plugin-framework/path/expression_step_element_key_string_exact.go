// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"fmt"
)

// Ensure ExpressionStepElementKeyStringExact satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepElementKeyStringExact("")

// ExpressionStepElementKeyStringExact is an attribute path expression for an exact string
// key within a map. Map keys are always strings.
type ExpressionStepElementKeyStringExact string

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyStringExact and the string element key is equivalent.
func (s ExpressionStepElementKeyStringExact) Equal(o ExpressionStep) bool {
	other, ok := o.(ExpressionStepElementKeyStringExact)

	if !ok {
		return false
	}

	return string(s) == string(other)
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyStringExact condition.
func (s ExpressionStepElementKeyStringExact) Matches(pathStep PathStep) bool {
	pathStepElementKeyString, ok := pathStep.(PathStepElementKeyString)

	if !ok {
		return false
	}

	return string(s) == string(pathStepElementKeyString)
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyStringExact) String() string {
	return fmt.Sprintf("[%q]", string(s))
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyStringExact) unexported() {}
