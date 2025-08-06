// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// BlockWithValidateImplementation is an optional interface on
// Block which enables validation of the provider-defined implementation
// for the Block. This logic runs during Validate* RPCs, or via
// provider-defined unit testing, to ensure the provider's definition is valid
// before further usage could cause other unexpected errors or panics.
type BlockWithValidateImplementation interface {
	Block

	// ValidateImplementation should contain the logic which validates
	// the Block implementation. Since this logic can prevent the provider
	// from being usable, it should be very targeted and defensive against
	// false positives.
	ValidateImplementation(context.Context, ValidateImplementationRequest, *ValidateImplementationResponse)
}

// ValidateBlockImplementation contains the generic Block implementation
// validation logic for all types.
//
// This logic currently:
//   - Checks whether the given AttributeName in the path is a valid identifier
//   - If the given Block implements the BlockWithValidateImplementation
//     interface, calls the method
//   - Recursively calls this function on nested attributes and blocks
func ValidateBlockImplementation(ctx context.Context, block Block, req ValidateImplementationRequest) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(IsReservedResourceAttributeName(req.Name, req.Path)...)
	diags.Append(IsValidAttributeName(req.Name, req.Path)...)

	if blockWithValidateImplementation, ok := block.(BlockWithValidateImplementation); ok {
		resp := &ValidateImplementationResponse{}

		blockWithValidateImplementation.ValidateImplementation(ctx, req, resp)

		diags.Append(resp.Diagnostics...)
	}

	nestedObject := block.GetNestedObject()

	if nestedObject == nil {
		return diags
	}

	nestingMode := block.GetNestingMode()

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
		// case BlockNestingModeList:
		// 	nestedAttributePath = req.Path.AtListIndex(0).AtName(nestedAttributeName)
		// case BlockNestingModeSet:
		// 	nestedAttributePath = req.Path.AtSetValue(types.StringValue("*")).AtName(nestedAttributeName)
		// case BlockNestingModeSingle:
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

	for nestedBlockName, nestedBlock := range nestedObject.GetBlocks() {
		var nestedBlockPath path.Path

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
		// case BlockNestingModeList:
		// 	nestedBlockPath = req.Path.AtListIndex(0).AtName(nestedBlockName)
		// case BlockNestingModeSet:
		// 	nestedBlockPath = req.Path.AtSetValue(types.StringValue("*")).AtName(nestedBlockName)
		// case BlockNestingModeSingle:
		// 	nestedBlockPath = req.Path.AtName(nestedBlockName)
		default:
			// This is purely to preserve the prior logic. Refer to above comment.
			nestedBlockPath = req.Path.AtName(nestedBlockName)
		}

		nestedReq := ValidateImplementationRequest{
			Name: nestedBlockName,
			Path: nestedBlockPath,
		}

		diags.Append(ValidateBlockImplementation(ctx, nestedBlock, nestedReq)...)
	}

	return diags
}
