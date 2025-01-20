// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Number is a schema validator for types.Number attributes.
type Number interface {
	Describer

	// ValidateNumber should perform the validation.
	ValidateNumber(context.Context, NumberRequest, *NumberResponse)
}

// NumberRequest is a request for types.Number schema validation.
type NumberRequest struct {
	// Path contains the path of the attribute for validation. Use this path
	// for any response diagnostics.
	Path path.Path

	// PathExpression contains the expression matching the exact path
	// of the attribute for validation.
	PathExpression path.Expression

	// Config contains the entire configuration of the data source, provider, or resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for validation from the configuration.
	ConfigValue types.Number
}

// NumberResponse is a response to a NumberRequest.
type NumberResponse struct {
	// Diagnostics report errors or warnings related to validating the data source, provider, or resource
	// configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
