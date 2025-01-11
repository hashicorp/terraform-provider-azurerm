// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.List = sizeAtLeastValidator{}
var _ function.ListParameterValidator = sizeAtLeastValidator{}

type sizeAtLeastValidator struct {
	min int
}

func (v sizeAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("list must contain at least %d elements", v.min)
}

func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sizeAtLeastValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) < v.min {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

func (v sizeAtLeastValidator) ValidateParameterList(ctx context.Context, req function.ListParameterValidatorRequest, resp *function.ListParameterValidatorResponse) {
	if req.Value.IsNull() || req.Value.IsUnknown() {
		return
	}

	elems := req.Value.Elements()

	if len(elems) < v.min {
		resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
			req.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		)
	}
}

// SizeAtLeast returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a List.
//   - Contains at least min elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtLeast(minVal int) sizeAtLeastValidator {
	return sizeAtLeastValidator{
		min: minVal,
	}
}
