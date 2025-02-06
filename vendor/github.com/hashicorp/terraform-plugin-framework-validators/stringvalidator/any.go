// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Any returns a validator which ensures that any configured attribute value
// passes at least one of the given validators.
//
// To prevent practitioner confusion should non-passing validators have
// conflicting logic, only warnings from the passing validator are returned.
// Use AnyWithAllWarnings() to return warnings from non-passing validators
// as well.
func Any(validators ...validator.String) validator.String {
	return anyValidator{
		validators: validators,
	}
}

var _ validator.String = anyValidator{}

// anyValidator implements the validator.
type anyValidator struct {
	validators []validator.String
}

// Description describes the validation in plain text formatting.
func (v anyValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v anyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	for _, subValidator := range v.validators {
		validateResp := &validator.StringResponse{}

		subValidator.ValidateString(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			resp.Diagnostics = validateResp.Diagnostics

			return
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
