// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Float64 = oneOfValidator{}
var _ function.Float64ParameterValidator = oneOfValidator{}

type oneOfValidator struct {
	values []types.Float64
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

func (v oneOfValidator) ValidateParameterFloat64(ctx context.Context, request function.Float64ParameterValidatorRequest, response *function.Float64ParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
		request.ArgumentPosition,
		v.Description(ctx),
		value.String(),
	)
}

// OneOf checks that the float64 held in the attribute or function parameter
// is one of the given `values`.
func OneOf(values ...float64) oneOfValidator {
	frameworkValues := make([]types.Float64, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.Float64Value(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}
