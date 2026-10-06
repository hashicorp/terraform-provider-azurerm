// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.String = lengthAtLeastValidator{}
var _ function.StringParameterValidator = lengthAtLeastValidator{}

type lengthAtLeastValidator struct {
	minLength int
}

func (validator lengthAtLeastValidator) invalidUsageMessage() string {
	return fmt.Sprintf("minLength cannot be less than zero - minLength: %d", validator.minLength)
}

func (validator lengthAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at least %d", validator.minLength)
}

func (validator lengthAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v lengthAtLeastValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 {
		response.Diagnostics.Append(
			validatordiag.InvalidValidatorUsageDiagnostic(
				request.Path,
				"LengthAtLeast",
				v.invalidUsageMessage(),
			),
		)

		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if l := len(value); l < v.minLength {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

func (v lengthAtLeastValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	// Return an error if the validator has been created in an invalid state
	if v.minLength < 0 {
		response.Error = validatorfuncerr.InvalidValidatorUsageFuncError(
			request.ArgumentPosition,
			"LengthAtLeast",
			v.invalidUsageMessage(),
		)

		return
	}

	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if l := len(value); l < v.minLength {
		response.Error = validatorfuncerr.InvalidParameterValueLengthFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", l),
		)

		return
	}
}

// LengthAtLeast returns an validator which ensures that any configured
// attribute or function parameter value is of single-byte character length greater than or equal
// to the given minimum. Null (unconfigured) and unknown (known after apply)
// values are skipped.
//
// minLength cannot be less than zero. Invalid input for minLength will result in an
// implementation error message during validation.
//
// Use UTF8LengthAtLeast for checking multiple-byte characters.
func LengthAtLeast(minLength int) lengthAtLeastValidator {
	return lengthAtLeastValidator{
		minLength: minLength,
	}
}
