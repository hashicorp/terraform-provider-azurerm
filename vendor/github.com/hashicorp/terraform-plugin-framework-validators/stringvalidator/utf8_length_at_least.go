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

var _ validator.String = utf8LengthAtLeastValidator{}
var _ function.StringParameterValidator = utf8LengthAtLeastValidator{}

type utf8LengthAtLeastValidator struct {
	minLength int
}

func (validator utf8LengthAtLeastValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minLength cannot be less than zero - minLength: %d", validator.minLength)
}

func (validator utf8LengthAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be at least %d", validator.minLength)
}

func (validator utf8LengthAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v utf8LengthAtLeastValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"UTF8LengthAtLeast",
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

	if count < v.minLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		))

		return
	}
}

func (v utf8LengthAtLeastValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"UTF8LengthAtLeast",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	count := utf8.RuneCountInString(value)

	if count < v.minLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", count),
		)

		return
	}
}

// UTF8LengthAtLeast returns an validator which ensures that any configured
// attribute or function parameter value is of UTF-8 character count greater than or equal to the
// given minimum. Null (unconfigured) and unknown (known after apply) values
// are skipped.
//
// minLength cannot be less than zero. Invalid input for minLength will result in an
// implementation error message during validation.
//
// Use LengthAtLeast for checking single-byte character counts.
func UTF8LengthAtLeast(minLength int) utf8LengthAtLeastValidator {
	return utf8LengthAtLeastValidator{
		minLength: minLength,
	}
}
