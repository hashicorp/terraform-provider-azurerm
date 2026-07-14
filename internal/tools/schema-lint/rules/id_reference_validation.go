// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var _ PropertyRule = idReferenceValidation{}

// idReferenceValidation flags user-settable "*_id" reference properties that
// have no validation, or only weak validation that accepts arbitrary strings
// (the validation.StringIsNotEmpty signature). Resource ID references should use
// a resource-specific ID validator (e.g. a commonids Validate...ID function) or
// validation.IsUUID.
type idReferenceValidation struct{}

// idProbes are clearly-invalid ID values that a resource-specific or UUID
// validator rejects. A validator that accepts all of them does not meaningfully
// constrain the ID.
var idProbes = []string{"x", "not-a-valid-id", "0000"}

func (idReferenceValidation) ID() string   { return "SL013" }
func (idReferenceValidation) Name() string { return "id-reference-validation" }

func (idReferenceValidation) Description() string {
	return `A user-settable "*_id" reference should use a resource-specific ID validator (or validation.IsUUID), not just StringIsNotEmpty.`
}

func (idReferenceValidation) DefaultSeverity() Severity { return SeverityWarning }

func (r idReferenceValidation) CheckProperty(ctx PropertyContext) []Finding {
	// Only user-settable string properties whose name ends in "_id".
	if !ctx.Schema.Optional && !ctx.Schema.Required {
		return nil
	}
	if ctx.Schema.Type != TypeString {
		return nil
	}
	if !strings.HasSuffix(ctx.Name, "_id") {
		return nil
	}

	accept := ctx.Schema.AcceptsValue
	if accept == nil {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("ID reference %q has no validation; use a resource-specific ID validator (e.g. a commonids Validate...ID function) or validation.IsUUID", ctx.Path))}
	}

	if acceptsArbitraryID(accept) {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("ID reference %q uses weak validation (accepts arbitrary strings); use a resource-specific ID validator or validation.IsUUID", ctx.Path))}
	}

	return nil
}

// acceptsArbitraryID reports whether the validator accepts every clearly-invalid
// probe value, i.e. it does not meaningfully constrain the ID.
func acceptsArbitraryID(accept func(string) bool) bool {
	for _, p := range idProbes {
		if !accept(p) {
			return false
		}
	}

	return true
}
