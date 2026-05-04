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

// ValueInt64sAre returns an validator which ensures that any configured
// Int64 values passes each Int64 validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValueInt64sAre(elementValidators ...validator.Int64) validator.List {
	return valueInt64sAreValidator{
		elementValidators: elementValidators,
	}
}

var _ validator.List = valueInt64sAreValidator{}

// valueInt64sAreValidator validates that each Int64 member validates against each of the value validators.
type valueInt64sAreValidator struct {
	elementValidators []validator.Int64
}

// Description describes the validation in plain text formatting.
func (v valueInt64sAreValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, elementValidator := range v.elementValidators {
		descriptions = append(descriptions, elementValidator.Description(ctx))
	}

	return fmt.Sprintf("element value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valueInt64sAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateInt64 performs the validation.
func (v valueInt64sAreValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, ok := req.ConfigValue.ElementType(ctx).(basetypes.Int64Typable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			"While performing schema-based validation, an unexpected error occurred. "+
				"The attribute declares a Int64 values validator, however its values do not implement types.Int64Type or the types.Int64Typable interface for custom Int64 types. "+
				"Use the appropriate values validator that matches the element type. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", req.Path.String())+
				fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx)),
		)

		return
	}

	for idx, element := range req.ConfigValue.Elements() {
		elementPath := req.Path.AtListIndex(idx)

		elementValuable, ok := element.(basetypes.Int64Valuable)

		// The check above should have prevented this, but raise an error
		// instead of a type assertion panic or skipping the element. Any issue
		// here likely indicates something wrong in the framework itself.
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Validator for Element Value",
				"While performing schema-based validation, an unexpected error occurred. "+
					"The attribute declares a Int64 values validator, however its values do not implement types.Int64Type or the types.Int64Typable interface for custom Int64 types. "+
					"This is likely an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Path: %s\n", req.Path.String())+
					fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx))+
					fmt.Sprintf("Element Value Type: %T\n", element),
			)

			return
		}

		elementValue, diags := elementValuable.ToInt64Value(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early if the new diagnostics indicate an issue since
		// it likely will be the same for all elements.
		if diags.HasError() {
			return
		}

		elementReq := validator.Int64Request{
			Path:           elementPath,
			PathExpression: elementPath.Expression(),
			ConfigValue:    elementValue,
			Config:         req.Config,
		}

		for _, elementValidator := range v.elementValidators {
			elementResp := &validator.Int64Response{}

			elementValidator.ValidateInt64(ctx, elementReq, elementResp)

			resp.Diagnostics.Append(elementResp.Diagnostics...)
		}
	}
}
