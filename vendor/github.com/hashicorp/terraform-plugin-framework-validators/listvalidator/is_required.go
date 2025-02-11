// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = isRequiredValidator{}

// isRequiredValidator validates that a list has a configuration value.
type isRequiredValidator struct{}

// Description describes the validation in plain text formatting.
func (v isRequiredValidator) Description(_ context.Context) string {
	return "must have a configuration value as the provider has marked it as required"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v isRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v isRequiredValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() {
		resp.Diagnostics.Append(validatordiag.InvalidBlockDiagnostic(
			req.Path,
			v.Description(ctx),
		))
	}
}

// IsRequired returns a validator which ensures that any configured list has a value (not null).
//
// This validator is equivalent to the `Required` field on attributes and is only
// practical for use with `schema.ListNestedBlock`
func IsRequired() validator.List {
	return isRequiredValidator{}
}
