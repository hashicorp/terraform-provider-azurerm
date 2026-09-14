// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Float64 = betweenValidator{}
var _ function.Float64ParameterValidator = betweenValidator{}

type betweenValidator struct {
	min, max float64
}

func (validator betweenValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minVal cannot be greater than maxVal - minVal: %f, maxVal: %f", validator.min, validator.max)
}

func (validator betweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %f and %f", validator.min, validator.max)
}

func (validator betweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v betweenValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	// Return an error if the validator has been created in an invalid state
	if v.min > v.max {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"Between",
				v.invalidUsageMessage(),
			),
		)

		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat64()

	if value < v.min || value > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%f", value),
		))
	}
}

func (v betweenValidator) ValidateParameterFloat64(ctx context.Context, request function.Float64ParameterValidatorRequest, response *function.Float64ParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.min > v.max {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"Between",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueFloat64()

	if value < v.min || value > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%f", value),
		)
	}
}

// Between returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 64-bit floating point.
//   - Is greater than or equal to the given minimum and less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
//
// minVal cannot be greater than maxVal. Invalid combinations of
// minVal and maxVal will result in an implementation error message during validation.
func Between(minVal, maxVal float64) betweenValidator {
	return betweenValidator{
		min: minVal,
		max: maxVal,
	}
}
