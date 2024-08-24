// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = lengthAtMostValidator{}

// lengthAtMostValidator validates that a string Attribute's length is at most a certain value.
type lengthAtMostValidator struct {
	maxLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at most %d", validator.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v lengthAtMostValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// LengthAtMost returns an validator which ensures that any configured
// attribute value is of single-byte character length less than or equal
// to the given maximum. Null (unconfigured) and unknown (known after apply)
// values are skipped.
//
// Use UTF8LengthAtMost for checking multiple-byte characters.
func LengthAtMost(maxLength int) validator.String {
	if maxLength < 0 {
		return nil
	}

	return lengthAtMostValidator{
		maxLength: maxLength,
	}
}
