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

var _ validator.String = utf8LengthAtMostValidator{}

// utf8LengthAtMostValidator implements the validator.
type utf8LengthAtMostValidator struct {
	maxLength int
}

// Description describes the validation in plain text formatting.
func (validator utf8LengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be at most %d", validator.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator utf8LengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v utf8LengthAtMostValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// UTF8LengthAtMost returns an validator which ensures that any configured
// attribute value is of UTF-8 character count less than or equal to the
// given maximum. Null (unconfigured) and unknown (known after apply) values
// are skipped.
//
// Use LengthAtMost for checking single-byte character counts.
func UTF8LengthAtMost(maxLength int) validator.String {
	if maxLength < 0 {
		return nil
	}

	return utf8LengthAtMostValidator{
		maxLength: maxLength,
	}
}
