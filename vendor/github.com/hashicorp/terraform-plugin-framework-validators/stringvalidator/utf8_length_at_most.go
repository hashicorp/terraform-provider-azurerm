// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.String = utf8LengthAtMostValidator{}
var _ function.StringParameterValidator = utf8LengthAtMostValidator{}

type utf8LengthAtMostValidator struct {
	maxLength int
}

func (validator utf8LengthAtMostValidator) invalidUsageMessage() string {
	return fmt.Sprintf("maxLength cannot be less than zero - maxLength: %d", validator.maxLength)
}

func (validator utf8LengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be at most %d", validator.maxLength)
}

func (validator utf8LengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v utf8LengthAtMostValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.maxLength < 0 {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"UTF8LengthAtMost",
				v.invalidUsageMessage(),
			),
		)

		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	count := utf8.RuneCountInString(value)

	if count > v.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		))

		return
	}
}

func (v utf8LengthAtMostValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.maxLength < 0 {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"UTF8LengthAtMost",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	count := utf8.RuneCountInString(value)

	if count > v.maxLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		)

		return
	}
}

// UTF8LengthAtMost returns an validator which ensures that any configured
// attribute or function parameter value is of UTF-8 character count less than or equal to the
// given maximum. Null (unconfigured) and unknown (known after apply) values
// are skipped.
//
// maxLength cannot be less than zero. Invalid input for maxLength will result in an
// implementation error message during validation.
//
// Use LengthAtMost for checking single-byte character counts.
func UTF8LengthAtMost(maxLength int) utf8LengthAtMostValidator {
	return utf8LengthAtMostValidator{
		maxLength: maxLength,
	}
}
