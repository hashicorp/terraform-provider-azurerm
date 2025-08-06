// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "fmt"

// Ensure PathStepElementKeyString satisfies the PathStep interface.
var _ PathStep = PathStepElementKeyString("")

// PathStepElementKeyString is an attribute path transversal for a string
// key of a map. Map keys are always strings.
//
// List elements must be transversed by PathStepElementKeyInt.
// Object attributes must be transversed by PathStepAttributeName.
// Set elements must be transversed by PathStepElementKeyValue.
type PathStepElementKeyString string

// Equal returns true if the given PathStep is a PathStepAttributeName and the
// attribute name is equivalent.
func (s PathStepElementKeyString) Equal(o PathStep) bool {
	other, ok := o.(PathStepElementKeyString)

	if !ok {
		return false
	}

	return string(s) == string(other)
}

// ExpressionStep returns the ExpressionStep for the PathStep.
func (s PathStepElementKeyString) ExpressionStep() ExpressionStep {
	return ExpressionStepElementKeyStringExact(s)
}

// String returns the human-readable representation of the element key.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s PathStepElementKeyString) String() string {
	return fmt.Sprintf("[%q]", string(s))
}

// unexported satisfies the PathStep interface.
func (s PathStepElementKeyString) unexported() {}
