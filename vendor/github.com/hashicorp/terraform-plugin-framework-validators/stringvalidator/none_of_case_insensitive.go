// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = noneOfCaseInsensitiveValidator{}

// noneOfCaseInsensitiveValidator validates that the value matches one of expected values.
type noneOfCaseInsensitiveValidator struct {
	values []types.String
}

func (v noneOfCaseInsensitiveValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfCaseInsensitiveValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %s", v.values)
}

func (v noneOfCaseInsensitiveValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if strings.EqualFold(value.ValueString(), otherValue.ValueString()) {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				request.Path,
				v.Description(ctx),
				value.String(),
			))

			return
		}
	}
}

// NoneOfCaseInsensitive checks that the String held in the attribute
// is none of the given `values`.
func NoneOfCaseInsensitive(values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return noneOfCaseInsensitiveValidator{
		values: frameworkValues,
	}
}
