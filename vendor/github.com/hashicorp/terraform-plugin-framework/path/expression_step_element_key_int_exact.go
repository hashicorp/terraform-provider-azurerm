// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"fmt"
)

// Ensure ExpressionStepElementKeyIntExact satisfies the ExpressionStep
// interface.
var _ ExpressionStep = ExpressionStepElementKeyIntExact(0)

// ExpressionStepElementKeyIntExact is an attribute path expression for an exact integer
// element key match within a list. List indexing starts at 0.
type ExpressionStepElementKeyIntExact int64

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyIntExact and the integer element key is equivalent.
func (s ExpressionStepElementKeyIntExact) Equal(o ExpressionStep) bool {
	other, ok := o.(ExpressionStepElementKeyIntExact)

	if !ok {
		return false
	}

	return int64(s) == int64(other)
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyIntExact condition.
func (s ExpressionStepElementKeyIntExact) Matches(pathStep PathStep) bool {
	pathStepElementKeyInt, ok := pathStep.(PathStepElementKeyInt)

	if !ok {
		return false
	}

	return int64(s) == int64(pathStepElementKeyInt)
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyIntExact) String() string {
	return fmt.Sprintf("[%d]", s)
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyIntExact) unexported() {}
