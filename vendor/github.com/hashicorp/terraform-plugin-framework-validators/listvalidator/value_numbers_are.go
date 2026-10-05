// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueNumbersAre returns an validator which ensures that any configured
// Number values passes each Number validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValueNumbersAre(elementValidators ...validator.Number) validator.List {
	return valueNumbersAreValidator{
		elementValidators: elementValidators,
	}
}

var _ validator.List = valueNumbersAreValidator{}

// valueNumbersAreValidator validates that each Number member validates against each of the value validators.
type valueNumbersAreValidator struct {
	elementValidators []validator.Number
}

// Description describes the validation in plain text formatting.
func (v valueNumbersAreValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, elementValidator := range v.elementValidators {
		descriptions = append(descriptions, elementValidator.Description(ctx))
	}

	return fmt.Sprintf("element value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valueNumbersAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateNumber performs the validation.
func (v valueNumbersAreValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, ok := req.ConfigValue.ElementType(ctx).(basetypes.NumberTypable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			"While performing schema-based validation, an unexpected error occurred. "+
				"The attribute declares a Number values validator, however its values do not implement types.NumberType or the types.NumberTypable interface for custom Number types. "+
				"Use the appropriate values validator that matches the element type. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", req.Path.String())+
				fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx)),
		)

		return
	}

	for idx, element := range req.ConfigValue.Elements() {
		elementPath := req.Path.AtListIndex(idx)

		elementValuable, ok := element.(basetypes.NumberValuable)

		// The check above should have prevented this, but raise an error
		// instead of a type assertion panic or skipping the element. Any issue
		// here likely indicates something wrong in the framework itself.
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Validator for Element Value",
				"While performing schema-based validation, an unexpected error occurred. "+
					"The attribute declares a Number values validator, however its values do not implement types.NumberType or the types.NumberTypable interface for custom Number types. "+
					"This is likely an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Path: %s\n", req.Path.String())+
					fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx))+
					fmt.Sprintf("Element Value Type: %T\n", element),
			)

			return
		}

		elementValue, diags := elementValuable.ToNumberValue(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early if the new diagnostics indicate an issue since
		// it likely will be the same for all elements.
		if diags.HasError() {
			return
		}

		elementReq := validator.NumberRequest{
			Path:           elementPath,
			PathExpression: elementPath.Expression(),
			ConfigValue:    elementValue,
			Config:         req.Config,
		}

		for _, elementValidator := range v.elementValidators {
			elementResp := &validator.NumberResponse{}

			elementValidator.ValidateNumber(ctx, elementReq, elementResp)

			resp.Diagnostics.Append(elementResp.Diagnostics...)
		}
	}
}
