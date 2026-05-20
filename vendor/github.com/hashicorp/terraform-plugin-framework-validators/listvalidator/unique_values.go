// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = uniqueValuesValidator{}
var _ function.ListParameterValidator = uniqueValuesValidator{}

type uniqueValuesValidator struct{}

func (v uniqueValuesValidator) Description(_ context.Context) string {
	return "all values must be unique"
}

func (v uniqueValuesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

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

func (v uniqueValuesValidator) ValidateParameterList(ctx context.Context, req function.ListParameterValidatorRequest, resp *function.ListParameterValidatorResponse) {
	if req.Value.IsNull() || req.Value.IsUnknown() {
		return
	}

	elements := req.Value.Elements()

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

			resp.Error = function.ConcatFuncErrors(
				resp.Error,
				function.NewArgumentFuncError(
					req.ArgumentPosition,
					fmt.Sprintf("Duplicate List Value: This attribute contains duplicate values of: %s", elementInner),
				),
			)
		}
	}
}

// UniqueValues returns a validator which ensures that any configured list
// only contains unique values. This is similar to using a set attribute type
// which inherently validates unique values, but with list ordering semantics.
// Null (unconfigured) and unknown (known after apply) values are skipped.
func UniqueValues() uniqueValuesValidator {
	return uniqueValuesValidator{}
}
