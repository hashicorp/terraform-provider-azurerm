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

// ValueInt32sAre returns an validator which ensures that any configured
// Int32 values passes each Int32 validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValueInt32sAre(elementValidators ...validator.Int32) validator.List {
	return valueInt32sAreValidator{
		elementValidators: elementValidators,
	}
}

var _ validator.List = valueInt32sAreValidator{}

// valueInt32sAreValidator validates that each Int32 member validates against each of the value validators.
type valueInt32sAreValidator struct {
	elementValidators []validator.Int32
}

// Description describes the validation in plain text formatting.
func (v valueInt32sAreValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, elementValidator := range v.elementValidators {
		descriptions = append(descriptions, elementValidator.Description(ctx))
	}

	return fmt.Sprintf("element value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valueInt32sAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateInt32 performs the validation.
func (v valueInt32sAreValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, ok := req.ConfigValue.ElementType(ctx).(basetypes.Int32Typable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			"While performing schema-based validation, an unexpected error occurred. "+
				"The attribute declares a Int32 values validator, however its values do not implement types.Int32Type or the types.Int32Typable interface for custom Int32 types. "+
				"Use the appropriate values validator that matches the element type. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", req.Path.String())+
				fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx)),
		)

		return
	}

	for idx, element := range req.ConfigValue.Elements() {
		elementPath := req.Path.AtListIndex(idx)

		elementValuable, ok := element.(basetypes.Int32Valuable)

		// The check above should have prevented this, but raise an error
		// instead of a type assertion panic or skipping the element. Any issue
		// here likely indicates something wrong in the framework itself.
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Validator for Element Value",
				"While performing schema-based validation, an unexpected error occurred. "+
					"The attribute declares a Int32 values validator, however its values do not implement types.Int32Type or the types.Int32Typable interface for custom Int32 types. "+
					"This is likely an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Path: %s\n", req.Path.String())+
					fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx))+
					fmt.Sprintf("Element Value Type: %T\n", element),
			)

			return
		}

		elementValue, diags := elementValuable.ToInt32Value(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early if the new diagnostics indicate an issue since
		// it likely will be the same for all elements.
		if diags.HasError() {
			return
		}

		elementReq := validator.Int32Request{
			Path:           elementPath,
			PathExpression: elementPath.Expression(),
			ConfigValue:    elementValue,
			Config:         req.Config,
		}

		for _, elementValidator := range v.elementValidators {
			elementResp := &validator.Int32Response{}

			elementValidator.ValidateInt32(ctx, elementReq, elementResp)

			resp.Diagnostics.Append(elementResp.Diagnostics...)
		}
	}
}
