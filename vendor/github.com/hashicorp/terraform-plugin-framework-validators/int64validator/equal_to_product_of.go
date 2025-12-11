// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Int64 = equalToProductOfValidator{}

// equalToProductOfValidator validates that an integer Attribute's value equals the product of one
// or more integer Attributes retrieved via the given path expressions.
type equalToProductOfValidator struct {
	attributesToMultiplyPathExpressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (av equalToProductOfValidator) Description(_ context.Context) string {
	var attributePaths []string
	for _, p := range av.attributesToMultiplyPathExpressions {
		attributePaths = append(attributePaths, p.String())
	}

	return fmt.Sprintf("value must be equal to the product of %s", strings.Join(attributePaths, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (av equalToProductOfValidator) MarkdownDescription(ctx context.Context) string {
	return av.Description(ctx)
}

// ValidateInt64 performs the validation.
func (av equalToProductOfValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	// Ensure input path expressions resolution against the current attribute
	expressions := request.PathExpression.MergeExpressions(av.attributesToMultiplyPathExpressions...)

	// Multiply the value of all the attributes involved, but only if they are all known.
	productOfAttribs := int64(1)
	for _, expression := range expressions {
		matchedPaths, diags := request.Config.PathMatches(ctx, expression)
		response.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(request.Path) {
				continue
			}

			// Get the value
			var matchedValue attr.Value
			diags := request.Config.GetAttribute(ctx, mp, &matchedValue)
			response.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			if matchedValue.IsUnknown() {
				return
			}

			if matchedValue.IsNull() {
				return
			}

			// We know there is a value, convert it to the expected type
			var attribToMultiply types.Int64
			diags = tfsdk.ValueAs(ctx, matchedValue, &attribToMultiply)
			response.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			productOfAttribs *= attribToMultiply.ValueInt64()
		}
	}

	if request.ConfigValue.ValueInt64() != productOfAttribs {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			av.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

// EqualToProductOf returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is equal to the product of the given attributes retrieved via the given path expression(s).
//
// Validation is skipped if any null (unconfigured) and/or unknown (known after apply) values are present.
func EqualToProductOf(attributesToMultiplyPathExpressions ...path.Expression) validator.Int64 {
	return equalToProductOfValidator{attributesToMultiplyPathExpressions}
}
