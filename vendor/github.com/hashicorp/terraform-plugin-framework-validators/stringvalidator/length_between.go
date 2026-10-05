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

var _ validator.String = lengthBetweenValidator{}
var _ function.StringParameterValidator = lengthBetweenValidator{}

type lengthBetweenValidator struct {
	minLength, maxLength int
}

func (validator lengthBetweenValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minLength cannot be less than zero or greater than maxLength - minLength: %d, maxLength: %d", validator.minLength, validator.maxLength)
}

func (validator lengthBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", validator.minLength, validator.maxLength)
}

func (validator lengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v lengthBetweenValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 || v.minLength > v.maxLength {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"LengthBetween",
				v.invalidUsageMessage(),
			),
		)

		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if l := len(value); l < v.minLength || l > v.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

func (v lengthBetweenValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 || v.minLength > v.maxLength {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"LengthBetween",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if l := len(value); l < v.minLength || l > v.maxLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		)

		return
	}
}

// LengthBetween returns a validator which ensures that any configured
// attribute or function parameter value is of single-byte character length greater than or equal
// to the given minimum and less than or equal to the given maximum. Null
// (unconfigured) and unknown (known after apply) values are skipped.
//
// minLength cannot be less than zero or greater than maxLength. Invalid combinations of
// minLength and maxLength will result in an implementation error message during validation.
//
// Use UTF8LengthBetween for checking multiple-byte characters.
func LengthBetween(minLength, maxLength int) lengthBetweenValidator {
	return lengthBetweenValidator{
		minLength: minLength,
		maxLength: maxLength,
	}
}
