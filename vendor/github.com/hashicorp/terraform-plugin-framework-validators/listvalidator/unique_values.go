// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = uniqueValuesValidator{}

// uniqueValuesValidator implements the validator.
type uniqueValuesValidator struct{}

// Description returns the plaintext description of the validator.
func (v uniqueValuesValidator) Description(_ context.Context) string {
	return "all values must be unique"
}

// MarkdownDescription returns the Markdown description of the validator.
func (v uniqueValuesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateList implements the validation logic.
func (v uniqueValuesValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for indexOuter, elementOuter := range elements {
		// Only evaluate known values for duplicates.
		if elementOuter.IsUnknown() {
			continue
		}

		for indexInner := indexOuter + 1; indexInner < len(elements); indexInner++ {
			elementInner := elements[indexInner]

			if elementInner.IsUnknown() {
				continue
			}

			if !elementInner.Equal(elementOuter) {
				continue
			}

			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Duplicate List Value",
				fmt.Sprintf("This attribute contains duplicate values of: %s", elementInner),
			)
		}
	}
}

// UniqueValues returns a validator which ensures that any configured list
// only contains unique values. This is similar to using a set attribute type
// which inherently validates unique values, but with list ordering semantics.
// Null (unconfigured) and unknown (known after apply) values are skipped.
func UniqueValues() validator.List {
	return uniqueValuesValidator{}
}
