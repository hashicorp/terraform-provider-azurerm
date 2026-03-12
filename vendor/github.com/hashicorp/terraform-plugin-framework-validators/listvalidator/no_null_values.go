// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = noNullValuesValidator{}
var _ function.ListParameterValidator = noNullValuesValidator{}

type noNullValuesValidator struct{}

func (v noNullValuesValidator) Description(_ context.Context) string {
	return "All values in the list must be configured"
}

func (v noNullValuesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v noNullValuesValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for _, e := range elements {
		// Only evaluate known values for null
		if e.IsUnknown() {
			continue
		}

		if e.IsNull() {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Null List Value",
				"This attribute contains a null value.",
			)
		}
	}
}

func (v noNullValuesValidator) ValidateParameterList(ctx context.Context, req function.ListParameterValidatorRequest, resp *function.ListParameterValidatorResponse) {
	if req.Value.IsNull() || req.Value.IsUnknown() {
		return
	}

	elements := req.Value.Elements()

	for _, e := range elements {
		// Only evaluate known values for null
		if e.IsUnknown() {
			continue
		}

		if e.IsNull() {
			resp.Error = function.ConcatFuncErrors(
				resp.Error,
				function.NewArgumentFuncError(
					req.ArgumentPosition,
					"Null List Value: This attribute contains a null value.",
				),
			)
		}
	}
}

// NoNullValues returns a validator which ensures that any configured list
// only contains non-null values.
func NoNullValues() noNullValuesValidator {
	return noNullValuesValidator{}
}
