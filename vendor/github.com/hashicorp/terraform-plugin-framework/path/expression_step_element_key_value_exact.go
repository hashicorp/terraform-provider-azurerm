// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Ensure ExpressionStepElementKeyValueExact satisfies the Step interface.
// var _ Step = ExpressionStepElementKeyValueExact(/* ... */)

// ExpressionStepElementKeyValueExact is an attribute path expression for an exact Value
// element within a set. Sets do not use integer-based indexing.
type ExpressionStepElementKeyValueExact struct {
	// Value is an interface, so it cannot be type aliased with methods.
	attr.Value
}

// Equal returns true if the given ExpressionStep is a
// ExpressionStepElementKeyValueExact and the Value element key is equivalent.
func (s ExpressionStepElementKeyValueExact) Equal(o ExpressionStep) bool {
	other, ok := o.(ExpressionStepElementKeyValueExact)

	if !ok {
		return false
	}

	return s.Value.Equal(other.Value)
}

// Matches returns true if the given PathStep is fulfilled by the
// ExpressionStepElementKeyValueExact condition.
func (s ExpressionStepElementKeyValueExact) Matches(pathStep PathStep) bool {
	pathStepElementKeyValue, ok := pathStep.(PathStepElementKeyValue)

	if !ok {
		return false
	}

	return s.Value.Equal(pathStepElementKeyValue.Value)
}

// String returns the human-readable representation of the element key
// expression. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
func (s ExpressionStepElementKeyValueExact) String() string {
	return fmt.Sprintf("[Value(%s)]", s.Value.String())
}

// unexported satisfies the Step interface.
func (s ExpressionStepElementKeyValueExact) unexported() {}
