// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = lengthAtLeastValidator{}

// stringLenAtLeastValidator validates that a string Attribute's length is at least a certain value.
type lengthAtLeastValidator struct {
	minLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at least %d", validator.minLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v lengthAtLeastValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// LengthAtLeast returns an validator which ensures that any configured
// attribute value is of single-byte character length greater than or equal
// to the given minimum. Null (unconfigured) and unknown (known after apply)
// values are skipped.
//
// Use UTF8LengthAtLeast for checking multiple-byte characters.
func LengthAtLeast(minLength int) validator.String {
	if minLength < 0 {
		return nil
	}

	return lengthAtLeastValidator{
		minLength: minLength,
	}
}
