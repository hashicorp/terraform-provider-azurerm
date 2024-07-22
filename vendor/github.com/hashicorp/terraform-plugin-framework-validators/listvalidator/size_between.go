// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = sizeBetweenValidator{}

// sizeBetweenValidator validates that list contains at least min elements
// and at most max elements.
type sizeBetweenValidator struct {
	min int
	max int
}

// Description describes the validation in plain text formatting.
func (v sizeBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("list must contain at least %d elements and at most %d elements", v.min, v.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeBetweenValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) < v.min || len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

// SizeBetween returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a List.
//   - Contains at least min elements and at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeBetween(min, max int) validator.List {
	return sizeBetweenValidator{
		min: min,
		max: max,
	}
}
