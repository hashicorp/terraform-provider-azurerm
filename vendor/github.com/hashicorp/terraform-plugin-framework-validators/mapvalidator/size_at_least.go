// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Map = sizeAtLeastValidator{}
var _ function.MapParameterValidator = sizeAtLeastValidator{}

type sizeAtLeastValidator struct {
	min int
}

func (v sizeAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("map must contain at least %d elements", v.min)
}

func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sizeAtLeastValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
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

func (v sizeAtLeastValidator) ValidateParameterMap(ctx context.Context, req function.MapParameterValidatorRequest, resp *function.MapParameterValidatorResponse) {
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
//   - Is a Map.
//   - Contains at least min elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtLeast(minVal int) sizeAtLeastValidator {
	return sizeAtLeastValidator{
		min: minVal,
	}
}
