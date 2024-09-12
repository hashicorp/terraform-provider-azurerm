// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = utf8LengthAtLeastValidator{}

// utf8LengthAtLeastValidator implements the validator.
type utf8LengthAtLeastValidator struct {
	minLength int
}

// Description describes the validation in plain text formatting.
func (validator utf8LengthAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be at least %d", validator.minLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator utf8LengthAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v utf8LengthAtLeastValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// UTF8LengthAtLeast returns an validator which ensures that any configured
// attribute value is of UTF-8 character count greater than or equal to the
// given minimum. Null (unconfigured) and unknown (known after apply) values
// are skipped.
//
// Use LengthAtLeast for checking single-byte character counts.
func UTF8LengthAtLeast(minLength int) validator.String {
	if minLength < 0 {
		return nil
	}

	return utf8LengthAtLeastValidator{
		minLength: minLength,
	}
}
