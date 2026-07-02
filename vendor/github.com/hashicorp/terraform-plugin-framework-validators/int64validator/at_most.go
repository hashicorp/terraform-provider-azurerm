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

var _ validator.Int64 = atMostValidator{}
var _ function.Int64ParameterValidator = atMostValidator{}

type atMostValidator struct {
	max int64
}

func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %d", validator.max)
}

func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v atMostValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueInt64() > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

func (v atMostValidator) ValidateParameterInt64(ctx context.Context, request function.Int64ParameterValidatorRequest, response *function.Int64ParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	if request.Value.ValueInt64() > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", request.Value.ValueInt64()),
		)
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(maxVal int64) atMostValidator {
	return atMostValidator{
		max: maxVal,
	}
}
