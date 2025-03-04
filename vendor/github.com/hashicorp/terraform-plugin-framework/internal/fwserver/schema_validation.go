// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateSchemaRequest repesents a request for validating a Schema.
type ValidateSchemaRequest struct {
	// Config contains the entire configuration of the data source, provider, or resource.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config
}

// ValidateSchemaResponse represents a response to a
// ValidateSchemaRequest.
type ValidateSchemaResponse struct {
	// Diagnostics report errors or warnings related to validating the schema.
	// An empty slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics
}

// SchemaValidate performs all Attribute and Block validation.
//
// TODO: Clean up this abstraction back into an internal Schema type method.
// The extra Schema parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func SchemaValidate(ctx context.Context, s fwschema.Schema, req ValidateSchemaRequest, resp *ValidateSchemaResponse) {
	for name, attribute := range s.GetAttributes() {

		attributeReq := ValidateAttributeRequest{
			AttributePath:           path.Root(name),
			AttributePathExpression: path.MatchRoot(name),
			Config:                  req.Config,
		}
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		attributeResp := &ValidateAttributeResponse{}

		AttributeValidate(ctx, attribute, attributeReq, attributeResp)

		resp.Diagnostics.Append(attributeResp.Diagnostics...)
	}

	for name, block := range s.GetBlocks() {
		attributeReq := ValidateAttributeRequest{
			AttributePath:           path.Root(name),
			AttributePathExpression: path.MatchRoot(name),
			Config:                  req.Config,
		}
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		attributeResp := &ValidateAttributeResponse{}

		BlockValidate(ctx, block, attributeReq, attributeResp)

		resp.Diagnostics.Append(attributeResp.Diagnostics...)
	}

	if s.GetDeprecationMessage() != "" {
		resp.Diagnostics.AddWarning(
			"Deprecated",
			s.GetDeprecationMessage(),
		)
	}
}
