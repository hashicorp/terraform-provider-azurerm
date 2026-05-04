// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Expression represents an attribute path with expression steps, which can
// represent zero, one, or more actual paths in schema data. This logic is
// either based on an absolute path starting at the root of the schema data,
// similar to Path, or a relative path which is intended to be merged with an
// existing absolute path.
//
// Use the MatchRoot() function to create an Expression for an absolute path
// with an initial AtName() step. Use the MatchRelative() function to create
// an Expression for a relative path, which will be merged with an existing
// absolute path.
//
// Similar to Path, Expression functionality has some overlapping method names
// and follows a builder pattern, which allows for chaining method calls to
// construct a full expression. The available traversal steps after Expression
// creation are:
//
//   - AtAnyListIndex(): Step into a list at any index
//   - AtAnyMapKey(): Step into a map at any key
//   - AtAnySetValue(): Step into a set at any attr.Value element
//   - AtListIndex(): Step into a list at a specific index
//   - AtMapKey(): Step into a map at a specific key
//   - AtName(): Step into an attribute or block with a specific name
//   - AtParent(): Step backwards one step
//   - AtSetValue(): Step into a set at a specific attr.Value element
//
// For example, to express any list element with a root list attribute named
// "some_attribute":
//
//	path.MatchRoot("some_attribute").AtAnyListIndex()
//
// An Expression is generally preferable over a Path in schema-defined
// functionality that is intended to accept paths as parameters, such as
// attribute validators and attribute plan modifiers, since it allows consumers
// to support relative paths. Use the Merge() or MergeExpressions() method to
// combine the current attribute path expression with those expression(s).
//
// To find Paths from an Expression in schema based data structures, such as
// tfsdk.Config, tfsdk.Plan, and tfsdk.State, use their PathMatches() method.
type Expression struct {
	// root stores whether an expression was intentionally created to start
	// from the root of the data. This is used with Merge to overwrite steps
	// instead of appending steps.
	root bool

	// steps is the transversals included with the expression. In general,
	// operations against the path should protect against modification of the
	// original.
	steps ExpressionSteps
}

// AtAnyListIndex returns a copied expression with a new list index step at the
// end. The returned path is safe to modify without affecting the original.
func (e Expression) AtAnyListIndex() Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyIntAny{})

	return copiedPath
}

// AtAnyMapKey returns a copied expression with a new map key step at the end.
// The returned path is safe to modify without affecting the original.
func (e Expression) AtAnyMapKey() Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyStringAny{})

	return copiedPath
}

// AtAnySetValue returns a copied expression with a new set value step at the
// end. The returned path is safe to modify without affecting the original.
func (e Expression) AtAnySetValue() Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyValueAny{})

	return copiedPath
}

// AtListIndex returns a copied expression with a new list index step at the
// end. The returned path is safe to modify without affecting the original.
func (e Expression) AtListIndex(index int) Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyIntExact(index))

	return copiedPath
}

// AtMapKey returns a copied expression with a new map key step at the end.
// The returned path is safe to modify without affecting the original.
func (e Expression) AtMapKey(key string) Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyStringExact(key))

	return copiedPath
}

// AtName returns a copied expression with a new attribute or block name step
// at the end. The returned path is safe to modify without affecting the
// original.
func (e Expression) AtName(name string) Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepAttributeNameExact(name))

	return copiedPath
}

// AtParent returns a copied expression with a new parent step at the end.
// The returned path is safe to modify without affecting the original.
func (e Expression) AtParent() Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepParent{})

	return copiedPath
}

// AtSetValue returns a copied expression with a new set value step at the end.
// The returned path is safe to modify without affecting the original.
func (e Expression) AtSetValue(value attr.Value) Expression {
	copiedPath := e.Copy()

	copiedPath.steps.Append(ExpressionStepElementKeyValueExact{Value: value})

	return copiedPath
}

// Copy returns a duplicate of the expression that is safe to modify without
// affecting the original.
func (e Expression) Copy() Expression {
	return Expression{
		root:  e.root,
		steps: e.Steps().Copy(),
	}
}

// Equal returns true if the given expression is exactly equivalent.
func (e Expression) Equal(o Expression) bool {
	if e.root != o.root {
		return false
	}

	if e.steps == nil && o.steps == nil {
		return true
	}

	if e.steps == nil {
		return false
	}

	if !e.steps.Equal(o.steps) {
		return false
	}

	return true
}

// Matches returns true if the given Path is valid for the Expression. Any
// relative expression steps, such as ExpressionStepParent, are automatically
// resolved before matching.
func (e Expression) Matches(path Path) bool {
	return e.steps.Matches(path.Steps())
}

// MatchesParent returns true if the given Path is a valid parent for the
// Expression. This is helpful for determining if a child Path would
// potentially match the full Expression during depth-first traversal. Any
// relative expression steps, such as ExpressionStepParent, are automatically
// resolved before matching.
func (e Expression) MatchesParent(path Path) bool {
	return e.steps.MatchesParent(path.Steps())
}

// Merge returns a copied expression either with the steps of the given
// expression added to the end of the existing steps, or overwriting the
// steps if the given expression was a root expression.
//
// Any merged expressions will preserve all expressions steps, such as
// ExpressionStepParent, for troubleshooting. Methods such as Matches() will
// automatically resolve the expression when using it. Call the Resolve()
// method explicitly if a resolved expression without any ExpressionStepParent
// is desired.
func (e Expression) Merge(other Expression) Expression {
	if other.root {
		return other.Copy()
	}

	copiedExpression := e.Copy()

	copiedExpression.steps.Append(other.steps...)

	return copiedExpression
}

// MergeExpressions returns collection of expressions that calls Merge() on
// the current expression with each of the others.
//
// If no Expression are given, then it will return a collection of expressions
// containing only the current expression.
func (e Expression) MergeExpressions(others ...Expression) Expressions {
	var result Expressions

	if len(others) == 0 {
		result.Append(e)

		return result
	}

	for _, other := range others {
		result.Append(e.Merge(other))
	}

	return result
}

// Resolve returns a copied expression with any relative steps, such as
// ExpressionStepParent, resolved. This is not necessary before calling methods
// such as Matches(), however it can be useful before returning the String()
// method so the path information is simplified.
//
// Returns an empty expression if any ExpressionStepParent attempt to go
// beyond the first element.
func (e Expression) Resolve() Expression {
	copiedExpression := e.Copy()

	copiedExpression.steps = copiedExpression.steps.Resolve()

	return copiedExpression
}

// Steps returns a copy of the underlying expression steps. Returns an empty
// collection of steps if expression is nil.
func (e Expression) Steps() ExpressionSteps {
	if len(e.steps) == 0 {
		return ExpressionSteps{}
	}

	return e.steps.Copy()
}

// String returns the human-readable representation of the path.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (e Expression) String() string {
	return e.steps.String()
}

// MatchRelative creates an empty attribute path expression that is intended
// to be combined with an existing attribute path expression. This allows
// creating a relative expression in nested schemas, using AtParent() to
// traverse up the path or other At methods to traverse further down.
func MatchRelative() Expression {
	return Expression{
		steps: ExpressionSteps{},
	}
}

// MatchRoot creates an attribute path expression starting with
// ExpressionStepAttributeNameExact.
func MatchRoot(rootAttributeName string) Expression {
	return Expression{
		root: true,
		steps: ExpressionSteps{
			ExpressionStepAttributeNameExact(rootAttributeName),
		},
	}
}
