// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "strings"

// ExpressionSteps represents an ordered collection of attribute path
// expressions.
type ExpressionSteps []ExpressionStep

// Append adds the given ExpressionSteps to the end of the previous ExpressionSteps and
// returns the combined result.
func (s *ExpressionSteps) Append(steps ...ExpressionStep) ExpressionSteps {
	if s == nil {
		return steps
	}

	*s = append(*s, steps...)

	return *s
}

// Copy returns a duplicate of the steps that is safe to modify without
// affecting the original. Returns nil if the original steps is nil.
func (s ExpressionSteps) Copy() ExpressionSteps {
	if s == nil {
		return nil
	}

	copiedExpressionSteps := make(ExpressionSteps, len(s))

	copy(copiedExpressionSteps, s)

	return copiedExpressionSteps
}

// Equal returns true if the given ExpressionSteps are equivalent.
func (s ExpressionSteps) Equal(o ExpressionSteps) bool {
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

// LastStep returns the final ExpressionStep and the remaining ExpressionSteps.
func (s ExpressionSteps) LastStep() (ExpressionStep, ExpressionSteps) {
	if len(s) == 0 {
		return nil, ExpressionSteps{}
	}

	if len(s) == 1 {
		return s[0], ExpressionSteps{}
	}

	return s[len(s)-1], s[:len(s)-1]
}

// Matches returns true if the given PathSteps match each ExpressionStep.
//
// Any ExpressionStepParent will automatically be resolved.
func (s ExpressionSteps) Matches(pathSteps PathSteps) bool {
	resolvedExpressionSteps := s.Resolve()

	// Empty expression should not match anything to prevent false positives.
	if len(resolvedExpressionSteps) == 0 {
		return false
	}

	if len(resolvedExpressionSteps) != len(pathSteps) {
		return false
	}

	for stepIndex, expressionStep := range resolvedExpressionSteps {
		if !expressionStep.Matches(pathSteps[stepIndex]) {
			return false
		}
	}

	return true
}

// MatchesParent returns true if the given PathSteps match each ExpressionStep
// until there are no more PathSteps. This is helpful for determining if the
// PathSteps would potentially match the full ExpressionSteps during
// depth-first traversal.
//
// Any ExpressionStepParent will automatically be resolved.
func (s ExpressionSteps) MatchesParent(pathSteps PathSteps) bool {
	resolvedExpressionSteps := s.Resolve()

	// Empty expression should not match anything to prevent false positives.
	// Ensure to not return false on an empty path since walking a path always
	// starts with no steps.
	if len(resolvedExpressionSteps) == 0 {
		return false
	}

	// Path steps deeper than or equal to the expression steps should not match
	// as a potential parent.
	if len(pathSteps) >= len(resolvedExpressionSteps) {
		return false
	}

	for stepIndex, pathStep := range pathSteps {
		if !resolvedExpressionSteps[stepIndex].Matches(pathStep) {
			return false
		}
	}

	return true
}

// NextStep returns the first ExpressionStep and the remaining ExpressionSteps.
func (s ExpressionSteps) NextStep() (ExpressionStep, ExpressionSteps) {
	if len(s) == 0 {
		return nil, s
	}

	return s[0], s[1:]
}

// Resolve returns a copy of ExpressionSteps without any ExpressionStepParent.
//
// Returns empty ExpressionSteps if any ExpressionStepParent attempt to go
// beyond the first element. Returns nil if the original steps is nil.
func (s ExpressionSteps) Resolve() ExpressionSteps {
	if s == nil {
		return nil
	}

	result := make(ExpressionSteps, 0, len(s))

	// This might not be the most efficient or prettiest algorithm, but it
	// works for now.
	for _, step := range s {
		_, ok := step.(ExpressionStepParent)

		if !ok {
			result.Append(step)

			continue
		}

		// Allow parent traversal up to the root, but not beyond.
		if len(result) == 0 {
			return ExpressionSteps{}
		}

		_, remaining := result.LastStep()

		if len(remaining) == 0 {
			result = ExpressionSteps{}

			continue
		}

		result = remaining
	}

	return result
}

// String returns the human-readable representation of the ExpressionSteps.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (s ExpressionSteps) String() string {
	var result strings.Builder

	for stepIndex, step := range s {
		if stepIndex != 0 {
			switch step.(type) {
			case ExpressionStepAttributeNameExact, ExpressionStepParent:
				result.WriteString(".")
			}
		}

		result.WriteString(step.String())
	}

	return result.String()
}
