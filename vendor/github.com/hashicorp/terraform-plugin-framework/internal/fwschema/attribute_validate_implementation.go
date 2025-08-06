// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// AttributeWithValidateImplementation is an optional interface on
// Attribute which enables validation of the provider-defined implementation
// for the Attribute. This logic runs during Validate* RPCs, or via
// provider-defined unit testing, to ensure the provider's definition is valid
// before further usage could cause other unexpected errors or panics.
type AttributeWithValidateImplementation interface {
	Attribute

	// ValidateImplementation should contain the logic which validates
	// the Attribute implementation. Since this logic can prevent the provider
	// from being usable, it should be very targeted and defensive against
	// false positives.
	ValidateImplementation(context.Context, ValidateImplementationRequest, *ValidateImplementationResponse)
}

// ValidateImplementation contains the generic Attribute
// implementation validation logic for all types.
//
// This logic currently:
//   - Checks whether the given AttributeName in the path is a valid identifier
//   - If the given Attribute implements the
//     AttributeWithValidateImplementation interface, calls the method
//   - If the given Attribute implements the NestedAttribute interface,
//     recursively calls this function on nested attributes
func ValidateAttributeImplementation(ctx context.Context, attribute Attribute, req ValidateImplementationRequest) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(IsValidAttributeName(req.Name, req.Path)...)

	if attributeWithValidateImplementation, ok := attribute.(AttributeWithValidateImplementation); ok {
		resp := &ValidateImplementationResponse{}

		attributeWithValidateImplementation.ValidateImplementation(ctx, req, resp)

		diags.Append(resp.Diagnostics...)
	}

	nestedAttribute, ok := attribute.(NestedAttribute)

	if !ok {
		return diags
	}

	nestedObject := nestedAttribute.GetNestedObject()

	if nestedObject == nil {
		return diags
	}

	nestingMode := nestedAttribute.GetNestingMode()

	for nestedAttributeName, nestedAttribute := range nestedObject.GetAttributes() {
		var nestedAttributePath path.Path

		// TODO: path.Path and path.PathExpression are intended to map onto
		// actual data implementations, however we need some representation
		// for schema paths without data. It may make sense to introduce an
		// internal "schema path" to simplify outputting specialized
		// strings for these types of diagnostics.
		//
		// The below choices of AtListIndex(0), etc. are arbitrary in this
		// situation.
		//
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/574
		switch nestingMode {
		// case NestingModeList:
		// 	nestedAttributePath = req.Path.AtListIndex(0).AtName(nestedAttributeName)
		// case NestingModeMap:
		// 	nestedAttributePath = req.Path.AtMapKey("*").AtName(nestedAttributeName)
		// case NestingModeSet:
		// 	nestedAttributePath = req.Path.AtSetValue(types.StringValue("*")).AtName(nestedAttributeName)
		// case NestingModeSingle:
		// 	nestedAttributePath = req.Path.AtName(nestedAttributeName)
		default:
			// This is purely to preserve the prior logic. Refer to above comment.
			nestedAttributePath = req.Path.AtName(nestedAttributeName)
		}

		nestedReq := ValidateImplementationRequest{
			Name: nestedAttributeName,
			Path: nestedAttributePath,
		}

		diags.Append(ValidateAttributeImplementation(ctx, nestedAttribute, nestedReq)...)
	}

	return diags
}
