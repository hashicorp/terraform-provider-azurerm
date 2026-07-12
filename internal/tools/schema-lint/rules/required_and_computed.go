// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = requiredAndComputed{}

// requiredAndComputed flags a property that is declared both Required and
// Computed, which is an invalid schema combination.
type requiredAndComputed struct{}

func (requiredAndComputed) ID() string   { return "SL003" }
func (requiredAndComputed) Name() string { return "required-and-computed" }

func (requiredAndComputed) Description() string {
	return "A property must not be both Required and Computed."
}

func (requiredAndComputed) DefaultSeverity() Severity { return SeverityError }

func (r requiredAndComputed) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.Required && ctx.Schema.Computed {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("property %q is both Required and Computed", ctx.Path))}
	}

	return nil
}
