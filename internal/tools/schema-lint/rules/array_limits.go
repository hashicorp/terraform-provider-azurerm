// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = arrayLimits{}

// arrayLimits flags a user-settable scalar array (a TypeList/TypeSet of
// primitives) that does not set MaxItems. Arrays should declare MinItems and
// MaxItems based on the API constraints.
type arrayLimits struct{}

func (arrayLimits) ID() string   { return "SL007" }
func (arrayLimits) Name() string { return "array-limits" }

func (arrayLimits) Description() string {
	return "A scalar array (TypeList/TypeSet of primitives) should declare MinItems and MaxItems based on the API constraints."
}

func (arrayLimits) DefaultSeverity() Severity { return SeverityWarning }

func (r arrayLimits) CheckProperty(ctx PropertyContext) []Finding {
	if !ctx.Schema.Optional && !ctx.Schema.Required {
		return nil
	}
	if !IsCollection(ctx.Schema) {
		return nil
	}
	// Only scalar arrays; nested blocks are out of scope.
	if _, isBlock := BlockElem(ctx.Schema); isBlock {
		return nil
	}
	// A missing Elem is malformed schema and out of scope here.
	if ctx.Schema.Elem == nil {
		return nil
	}
	if ctx.Schema.MaxItems != 0 {
		return nil
	}

	return []Finding{propertyFinding(r, ctx, fmt.Sprintf("array property %q does not set MaxItems; declare MinItems/MaxItems based on the API constraints", ctx.Path))}
}
