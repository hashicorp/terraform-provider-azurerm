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

var _ validator.String = oneOfCaseInsensitiveValidator{}

// oneOfCaseInsensitiveValidator validates that the value matches one of expected values.
type oneOfCaseInsensitiveValidator struct {
	values []types.String
}

func (v oneOfCaseInsensitiveValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfCaseInsensitiveValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %s", v.values)
}

func (v oneOfCaseInsensitiveValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if strings.EqualFold(value.ValueString(), otherValue.ValueString()) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

// OneOfCaseInsensitive checks that the String held in the attribute
// is one of the given `values`.
func OneOfCaseInsensitive(values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return oneOfCaseInsensitiveValidator{
		values: frameworkValues,
	}
}
