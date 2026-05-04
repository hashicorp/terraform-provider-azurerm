// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.String = noneOfValidator{}
var _ function.StringParameterValidator = noneOfValidator{}

type noneOfValidator struct {
	values []types.String
}

func (v noneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %s", v.values)
}

func (v noneOfValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if !value.Equal(otherValue) {
			continue
		}

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))

		break
	}
}

func (v noneOfValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value

	for _, otherValue := range v.values {
		if !value.Equal(otherValue) {
			continue
		}

		response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			value.String(),
		)

		break
	}
}

// NoneOf checks that the String held in the attribute or function parameter
// is none of the given `values`.
func NoneOf(values ...string) noneOfValidator {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return noneOfValidator{
		values: frameworkValues,
	}
}
