// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// ExpressionStep represents an expression of an attribute path step, which may
// match zero, one, or more actual paths.
type ExpressionStep interface {
	// Equal should return true if the given Step is exactly equivalent.
	Equal(ExpressionStep) bool

	// Matches should return true if the given PathStep can be fulfilled by the
	// ExpressionStep.
	Matches(PathStep) bool

	// String should return a human-readable representation of the step
	// intended for logging and error messages. There should not be usage
	// that needs to be protected by compatibility guarantees.
	String() string

	// unexported prevents outside types from satisfying the interface.
	unexported()
}
