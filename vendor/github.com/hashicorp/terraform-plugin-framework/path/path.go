// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Path represents exact traversal steps into a schema or schema-based data.
// These steps always start from the root of the schema, which is an object
// with zero or more attributes and blocks.
//
// Use the Root() function to create a Path with an initial AtName() step. Path
// functionality follows a builder pattern, which allows for chaining method
// calls to construct a full path. The available traversal steps after Path
// creation are:
//
//   - AtListIndex(): Step into a list at a specific 0-based index
//   - AtMapKey(): Step into a map at a specific key
//   - AtName(): Step into an attribute or block with a specific name
//   - AtSetValue(): Step into a set at a specific attr.Value element
//
// For example, to represent the first list element with a root list attribute
// named "some_attribute":
//
//	path.MatchRoot("some_attribute").AtListIndex(0)
//
// Path is used for functionality which must exactly match the underlying
// schema structure and types, such as diagnostics that are intended for a
// specific attribute or working with specific attribute values in a schema
// based data structure such as tfsdk.Config, tfsdk.Plan, or tfsdk.State.
//
// Refer to Expression for situations where relative or wildcard step logic is
// desirable for schema defined functionality, such as attribute validators or
// attribute plan modifiers.
type Path struct {
	// steps is the transversals included with the path. In general, operations
	// against the path should protect against modification of the original.
	steps PathSteps
}

// AtListIndex returns a copied path with a new list index step at the end.
// The returned path is safe to modify without affecting the original.
//
// List indices are 0-based. The first element of a list is 0.
func (p Path) AtListIndex(index int) Path {
	copiedPath := p.Copy()

	copiedPath.steps.Append(PathStepElementKeyInt(index))

	return copiedPath
}

// AtTupleIndex returns a copied path with a new tuple index step at the end.
// The returned path is safe to modify without affecting the original.
//
// Tuple indices are 0-based. The first element of a tuple is 0.
func (p Path) AtTupleIndex(index int) Path {
	copiedPath := p.Copy()

	copiedPath.steps.Append(PathStepElementKeyInt(index))

	return copiedPath
}

// AtMapKey returns a copied path with a new map key step at the end.
// The returned path is safe to modify without affecting the original.
func (p Path) AtMapKey(key string) Path {
	copiedPath := p.Copy()

	copiedPath.steps.Append(PathStepElementKeyString(key))

	return copiedPath
}

// AtName returns a copied path with a new attribute or block name step at the
// end. The returned path is safe to modify without affecting the original.
func (p Path) AtName(name string) Path {
	copiedPath := p.Copy()

	copiedPath.steps.Append(PathStepAttributeName(name))

	return copiedPath
}

// AtSetValue returns a copied path with a new set value step at the end.
// The returned path is safe to modify without affecting the original.
func (p Path) AtSetValue(value attr.Value) Path {
	copiedPath := p.Copy()

	copiedPath.steps.Append(PathStepElementKeyValue{Value: value})

	return copiedPath
}

// Copy returns a duplicate of the path that is safe to modify without
// affecting the original.
func (p Path) Copy() Path {
	return Path{
		steps: p.Steps(),
	}
}

// Equal returns true if the given path is exactly equivalent.
func (p Path) Equal(o Path) bool {
	if p.steps == nil && o.steps == nil {
		return true
	}

	if p.steps == nil {
		return false
	}

	if !p.steps.Equal(o.steps) {
		return false
	}

	return true
}

// Expression returns an Expression which exactly matches the Path.
func (p Path) Expression() Expression {
	return Expression{
		root:  true,
		steps: p.steps.ExpressionSteps(),
	}
}

// ParentPath returns a copy of the path with the last step removed.
//
// If the current path is empty, an empty path is returned.
func (p Path) ParentPath() Path {
	if len(p.steps) == 0 {
		return Empty()
	}

	_, remainingSteps := p.steps.Copy().LastStep()

	return Path{
		steps: remainingSteps,
	}
}

// Steps returns a copy of the underlying path steps. Returns an empty
// collection of steps if path is nil.
func (p Path) Steps() PathSteps {
	if len(p.steps) == 0 {
		return PathSteps{}
	}

	return p.steps.Copy()
}

// String returns the human-readable representation of the path.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
func (p Path) String() string {
	return p.steps.String()
}

// Empty creates an empty attribute path. Provider code should use Root.
func Empty() Path {
	return Path{
		steps: PathSteps{},
	}
}

// Root creates an attribute path starting with a PathStepAttributeName.
func Root(rootAttributeName string) Path {
	return Path{
		steps: PathSteps{
			PathStepAttributeName(rootAttributeName),
		},
	}
}
