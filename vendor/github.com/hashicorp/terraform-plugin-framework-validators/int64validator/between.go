// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Int64 = betweenValidator{}
var _ function.Int64ParameterValidator = betweenValidator{}

type betweenValidator struct {
	min, max int64
}

func (validator betweenValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minVal cannot be greater than maxVal - minVal: %d, maxVal: %d", validator.min, validator.max)
}

func (validator betweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %d and %d", validator.min, validator.max)
}

func (validator betweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v betweenValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
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

	if request.ConfigValue.ValueInt64() < v.min || request.ConfigValue.ValueInt64() > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

func (v betweenValidator) ValidateParameterInt64(ctx context.Context, request function.Int64ParameterValidatorRequest, response *function.Int64ParameterValidatorResponse) {
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

	if request.Value.ValueInt64() < v.min || request.Value.ValueInt64() > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", request.Value.ValueInt64()),
		)
	}
}

// Between returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is greater than or equal to the given minimum and less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
//
// minVal cannot be greater than maxVal. Invalid combinations of
// minVal and maxVal will result in an implementation error message during validation.
func Between(minVal, maxVal int64) betweenValidator {
	return betweenValidator{
		min: minVal,
		max: maxVal,
	}
}
