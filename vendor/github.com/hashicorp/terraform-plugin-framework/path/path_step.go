// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

// PathStep represents a transversal for an attribute path. Only exact path
// transversals are supported as implementations of this interface must remain
// compatible with all protocol implementations.
type PathStep interface {
	// Equal should return true if the given PathStep is exactly equivalent.
	Equal(PathStep) bool

	// ExpressionStep should return an ExpressionStep which exactly
	// matches the PathStep.
	ExpressionStep() ExpressionStep

	// String should return a human-readable representation of the step
	// intended for logging and error messages. There should not be usage
	// that needs to be protected by compatibility guarantees.
	String() string

	// unexported prevents outside types from satisfying the interface.
	unexported()
}
