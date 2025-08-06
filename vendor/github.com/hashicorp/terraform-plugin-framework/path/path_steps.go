// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "strings"

// PathSteps represents an ordered collection of attribute path transversals.
type PathSteps []PathStep

// Append adds the given PathSteps to the end of the previous PathSteps and
// returns the combined result.
func (s *PathSteps) Append(steps ...PathStep) PathSteps {
	if s == nil {
		return steps
	}

	*s = append(*s, steps...)

	return *s
}

// Copy returns a duplicate of the steps that is safe to modify without
// affecting the original. Returns nil if the original steps is nil.
func (s PathSteps) Copy() PathSteps {
	if s == nil {
		return nil
	}

	copiedPathSteps := make(PathSteps, len(s))

	copy(copiedPathSteps, s)

	return copiedPathSteps
}

// Equal returns true if the given PathSteps are equivalent.
func (s PathSteps) Equal(o PathSteps) bool {
	if len(s) != len(o) {
		return false
	}

	for stepIndex, step := range s {
		if !step.Equal(o[stepIndex]) {
			return false
		}
	}

	return true
}

// LastStep returns the final PathStep and the remaining PathSteps.
func (s PathSteps) LastStep() (PathStep, PathSteps) {
	if len(s) == 0 {
		return nil, PathSteps{}
	}

	if len(s) == 1 {
		return s[0], PathSteps{}
	}

	return s[len(s)-1], s[:len(s)-1]
}

// NextStep returns the first PathStep and the remaining PathSteps.
func (s PathSteps) NextStep() (PathStep, PathSteps) {
	if len(s) == 0 {
		return nil, s
	}

	return s[0], s[1:]
}

// String returns the human-readable representation of the PathSteps.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s PathSteps) String() string {
	var result strings.Builder

	for stepIndex, step := range s {
		if _, ok := step.(PathStepAttributeName); ok && stepIndex != 0 {
			result.WriteString(".")
		}

		result.WriteString(step.String())
	}

	return result.String()
}

// ExpressionSteps returns the ordered collection of expression steps which
// exactly matches the PathSteps.
func (s PathSteps) ExpressionSteps() ExpressionSteps {
	result := make(ExpressionSteps, len(s))

	for stepIndex, pathStep := range s {
		result[stepIndex] = pathStep.ExpressionStep()
	}

	return result
}
