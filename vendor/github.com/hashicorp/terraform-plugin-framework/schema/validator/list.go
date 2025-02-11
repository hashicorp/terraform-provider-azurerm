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

// List is a schema validator for types.List attributes.
type List interface {
	Describer

	// ValidateList should perform the validation.
	ValidateList(context.Context, ListRequest, *ListResponse)
}

// ListRequest is a request for types.List schema validation.
type ListRequest struct {
	// Path contains the path of the attribute for validation. Use this path
	// for any response diagnostics.
	Path path.Path

	// PathExpression contains the expression matching the exact path
	// of the attribute for validation.
	PathExpression path.Expression

	// Config contains the entire configuration of the data source, provider, or resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for validation from the configuration.
	ConfigValue types.List
}

// ListResponse is a response to a ListRequest.
type ListResponse struct {
	// Diagnostics report errors or warnings related to validating the data source, provider, or resource
	// configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
