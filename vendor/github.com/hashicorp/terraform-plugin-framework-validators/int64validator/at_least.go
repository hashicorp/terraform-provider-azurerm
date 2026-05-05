// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Int64 = atLeastValidator{}
var _ function.Int64ParameterValidator = atLeastValidator{}

type atLeastValidator struct {
	min int64
}

func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %d", validator.min)
}

func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v atLeastValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueInt64() < v.min {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

func (v atLeastValidator) ValidateParameterInt64(ctx context.Context, request function.Int64ParameterValidatorRequest, response *function.Int64ParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	if request.Value.ValueInt64() < v.min {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", request.Value.ValueInt64()),
		)
	}
}

// AtLeast returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is greater than or equal to the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeast(minVal int64) atLeastValidator {
	return atLeastValidator{
		min: minVal,
	}
}
