// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueFloat32sAre returns an validator which ensures that any configured
// Float32 values passes each Float32 validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValueFloat32sAre(elementValidators ...validator.Float32) validator.Map {
	return valueFloat32sAreValidator{
		elementValidators: elementValidators,
	}
}

var _ validator.Map = valueFloat32sAreValidator{}

// valueFloat32sAreValidator validates that each Float32 member validates against each of the value validators.
type valueFloat32sAreValidator struct {
	elementValidators []validator.Float32
}

// Description describes the validation in plain text formatting.
func (v valueFloat32sAreValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, elementValidator := range v.elementValidators {
		descriptions = append(descriptions, elementValidator.Description(ctx))
	}

	return fmt.Sprintf("element value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valueFloat32sAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateFloat32 performs the validation.
func (v valueFloat32sAreValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, ok := req.ConfigValue.ElementType(ctx).(basetypes.Float32Typable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			"While performing schema-based validation, an unexpected error occurred. "+
				"The attribute declares a Float32 values validator, however its values do not implement types.Float32Type or the types.Float32Typable interface for custom Float32 types. "+
				"Use the appropriate values validator that matches the element type. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", req.Path.String())+
				fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx)),
		)

		return
	}

	for key, element := range req.ConfigValue.Elements() {
		elementPath := req.Path.AtMapKey(key)

		elementValuable, ok := element.(basetypes.Float32Valuable)

		// The check above should have prevented this, but raise an error
		// instead of a type assertion panic or skipping the element. Any issue
		// here likely indicates something wrong in the framework itself.
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Validator for Element Value",
				"While performing schema-based validation, an unexpected error occurred. "+
					"The attribute declares a Float32 values validator, however its values do not implement types.Float32Type or the types.Float32Typable interface for custom Float32 types. "+
					"This is likely an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Path: %s\n", req.Path.String())+
					fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx))+
					fmt.Sprintf("Element Value Type: %T\n", element),
			)

			return
		}

		elementValue, diags := elementValuable.ToFloat32Value(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early if the new diagnostics indicate an issue since
		// it likely will be the same for all elements.
		if diags.HasError() {
			return
		}

		elementReq := validator.Float32Request{
			Path:           elementPath,
			PathExpression: elementPath.Expression(),
			ConfigValue:    elementValue,
			Config:         req.Config,
		}

		for _, elementValidator := range v.elementValidators {
			elementResp := &validator.Float32Response{}

			elementValidator.ValidateFloat32(ctx, elementReq, elementResp)

			resp.Diagnostics.Append(elementResp.Diagnostics...)
		}
	}
}
