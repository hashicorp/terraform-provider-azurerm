// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Map = noNullValuesValidator{}
var _ function.MapParameterValidator = noNullValuesValidator{}

type noNullValuesValidator struct{}

func (v noNullValuesValidator) Description(_ context.Context) string {
	return "All values in the map must be configured"
}

func (v noNullValuesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v noNullValuesValidator) ValidateMap(_ context.Context, req validator.MapRequest, resp *validator.MapResponse) {
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
				"Null Map Value",
				"This attribute contains a null value.",
			)
		}
	}
}

func (v noNullValuesValidator) ValidateParameterMap(ctx context.Context, req function.MapParameterValidatorRequest, resp *function.MapParameterValidatorResponse) {
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
					"Null Map Value: This attribute contains a null value.",
				),
			)
		}
	}
}

// NoNullValues returns a validator which ensures that any configured map
// only contains non-null values.
func NoNullValues() noNullValuesValidator {
	return noNullValuesValidator{}
}
