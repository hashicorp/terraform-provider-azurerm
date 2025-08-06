// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = lengthAtMostValidator{}
var _ function.StringParameterValidator = lengthAtMostValidator{}

type lengthAtMostValidator struct {
	maxLength int
}

func (validator lengthAtMostValidator) invalidUsageMessage() string {
	return fmt.Sprintf("maxLength cannot be less than zero - maxLength: %d", validator.maxLength)
}

func (validator lengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at most %d", validator.maxLength)
}

func (validator lengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v lengthAtMostValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.maxLength < 0 {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"LengthAtMost",
				v.invalidUsageMessage(),
			),
		)

		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if l := len(value); l > v.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

func (v lengthAtMostValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.maxLength < 0 {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"LengthAtMost",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if l := len(value); l > v.maxLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		)

		return
	}
}

// LengthAtMost returns an validator which ensures that any configured
// attribute or function parameter value is of single-byte character length less than or equal
// to the given maximum. Null (unconfigured) and unknown (known after apply)
// values are skipped.
//
// maxLength cannot be less than zero. Invalid input for maxLength will result in an
// implementation error message during validation.
//
// Use UTF8LengthAtMost for checking multiple-byte characters.
func LengthAtMost(maxLength int) lengthAtMostValidator {
	return lengthAtMostValidator{
		maxLength: maxLength,
	}
}
