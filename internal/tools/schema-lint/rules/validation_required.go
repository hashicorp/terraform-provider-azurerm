// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = validationRequired{}

// validationRequired flags a user-settable scalar argument (string, int or
// float) that does not set any validation.
type validationRequired struct{}

func (validationRequired) ID() string   { return "SL009" }
func (validationRequired) Name() string { return "validation-required" }

func (validationRequired) Description() string {
	return "User-settable string and numeric arguments should set validation (strings at least StringIsNotEmpty; numerics a valid range)."
}

func (validationRequired) DefaultSeverity() Severity { return SeverityWarning }

func (r validationRequired) CheckProperty(ctx PropertyContext) []Finding {
	if !ctx.Schema.Optional && !ctx.Schema.Required {
		return nil
	}
	if ctx.Schema.HasValidateFunc {
		return nil
	}

	switch ctx.Schema.Type {
	case TypeString:
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("string property %q has no validation; add a ValidateFunc (StringIsNotEmpty at minimum)", ctx.Path))}
	case TypeInt, TypeFloat:
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("numeric property %q has no validation; specify a valid range at minimum", ctx.Path))}
	}

	return nil
}
