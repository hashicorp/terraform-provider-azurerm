// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = propertyDescriptionRequired{}

// propertyDescriptionRequired ensures every schema property documents itself
// with a non-empty Description.
type propertyDescriptionRequired struct{}

func (propertyDescriptionRequired) ID() string   { return "SL001" }
func (propertyDescriptionRequired) Name() string { return "property-description-required" }

func (propertyDescriptionRequired) Description() string {
	return "Every resource and data source property should set a non-empty Description."
}

func (propertyDescriptionRequired) DefaultSeverity() Severity { return SeverityWarning }

func (r propertyDescriptionRequired) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.Description == "" {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("property %q is missing a description", ctx.Path))}
	}

	return nil
}
