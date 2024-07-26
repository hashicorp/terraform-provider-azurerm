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

var _ validator.String = utf8LengthBetweenValidator{}

// utf8LengthBetweenValidator implements the validator.
type utf8LengthBetweenValidator struct {
	maxLength int
	minLength int
}

// Description describes the validation in plain text formatting.
func (v utf8LengthBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("UTF-8 character count must be between %d and %d", v.minLength, v.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v utf8LengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v utf8LengthBetweenValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// UTF8LengthBetween returns an validator which ensures that any configured
// attribute value is of UTF-8 character count greater than or equal to the
// given minimum and less than or equal to the given maximum. Null
// (unconfigured) and unknown (known after apply) values are skipped.
//
// Use LengthBetween for checking single-byte character counts.
func UTF8LengthBetween(minLength int, maxLength int) validator.String {
	if minLength < 0 || maxLength < 0 || minLength > maxLength {
		return nil
	}

	return utf8LengthBetweenValidator{
		maxLength: maxLength,
		minLength: minLength,
	}
}
