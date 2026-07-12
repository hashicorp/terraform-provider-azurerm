// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = optionalAndRequired{}

// optionalAndRequired flags a property that is declared both Optional and
// Required, which is an invalid schema combination.
type optionalAndRequired struct{}

func (optionalAndRequired) ID() string   { return "SL002" }
func (optionalAndRequired) Name() string { return "optional-and-required" }

func (optionalAndRequired) Description() string {
	return "A property must not be both Optional and Required."
}

func (optionalAndRequired) DefaultSeverity() Severity { return SeverityError }

func (r optionalAndRequired) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.Optional && ctx.Schema.Required {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("property %q is both Optional and Required", ctx.Path))}
	}

	return nil
}
