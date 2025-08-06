// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "strings"

// Expressions is a collection of attribute path expressions.
//
// Refer to the Expression documentation for more details about intended usage.
type Expressions []Expression

// Append adds the given Expressions to the collection without duplication and
// returns the combined result.
func (e *Expressions) Append(expressions ...Expression) Expressions {
	if e == nil {
		return expressions
	}

	for _, newExpression := range expressions {
		if e.Contains(newExpression) {
			continue
		}

		*e = append(*e, newExpression)
	}

	return *e
}

// Contains returns true if the collection of expressions includes the given
// expression.
func (e Expressions) Contains(checkExpression Expression) bool {
	for _, expression := range e {
		if expression.Equal(checkExpression) {
			return true
		}
	}

	return false
}

// Matches returns true if one of the expressions in the collection matches the
// given path.
func (e Expressions) Matches(checkPath Path) bool {
	for _, expression := range e {
		if expression.Matches(checkPath) {
			return true
		}
	}

	return false
}

// String returns the human-readable representation of the expression
// collection. It is intended for logging and error messages and is not
// protected by compatibility guarantees.
//
// Empty expressions are skipped.
func (p Expressions) String() string {
	var result strings.Builder

	result.WriteString("[")

	for pathIndex, path := range p {
		if path.Equal(Expression{}) {
			continue
		}

		if pathIndex != 0 {
			result.WriteString(",")
		}

		result.WriteString(path.String())
	}

	result.WriteString("]")

	return result.String()
}
