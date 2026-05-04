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

var _ validator.String = utf8LengthBetweenValidator{}
var _ function.StringParameterValidator = utf8LengthBetweenValidator{}

type utf8LengthBetweenValidator struct {
	maxLength int
	minLength int
}

func (v utf8LengthBetweenValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minLength and maxLength cannot be less than zero and maxLength must be greater than or equal to minLength - minLength: %d, maxLength: %d", v.minLength, v.maxLength)
}

func (v utf8LengthBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be between %d and %d", v.minLength, v.maxLength)
}

func (v utf8LengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v utf8LengthBetweenValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 || v.maxLength < 0 || v.minLength > v.maxLength {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"UTF8LengthBetween",
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

	if count < v.minLength || count > v.maxLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		))

		return
	}
}

func (v utf8LengthBetweenValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 || v.maxLength < 0 || v.minLength > v.maxLength {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"UTF8LengthBetween",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	count := utf8.RuneCountInString(value)

	if count < v.minLength || count > v.maxLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		)

		return
	}
}

// UTF8LengthBetween returns an validator which ensures that any configured
// attribute or function parameter value is of UTF-8 character count greater than or equal to the
// given minimum and less than or equal to the given maximum. Null
// (unconfigured) and unknown (known after apply) values are skipped.
//
// minLength and maxLength cannot be less than zero and maxLength must be greater than or equal to minLength.
// Invalid combinations of minLength and maxLength will result in an implementation error message
// during validation.
//
// Use LengthBetween for checking single-byte character counts.
func UTF8LengthBetween(minLength int, maxLength int) utf8LengthBetweenValidator {
	return utf8LengthBetweenValidator{
		maxLength: maxLength,
		minLength: minLength,
	}
}
