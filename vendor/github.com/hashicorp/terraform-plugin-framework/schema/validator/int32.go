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

// Int32 is a schema validator for types.Int32 attributes.
type Int32 interface {
	Describer

	// ValidateInt32 should perform the validation.
	ValidateInt32(context.Context, Int32Request, *Int32Response)
}

// Int32Request is a request for types.Int32 schema validation.
type Int32Request struct {
	// Path contains the path of the attribute for validation. Use this path
	// for any response diagnostics.
	Path path.Path

	// PathExpression contains the expression matching the exact path
	// of the attribute for validation.
	PathExpression path.Expression

	// Config contains the entire configuration of the data source, provider, or resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for validation from the configuration.
	ConfigValue types.Int32

	// ClientCapabilities defines optionally supported protocol features for
	// schema validation RPCs, such as forward-compatible Terraform
	// behavior changes.
	ClientCapabilities ValidateSchemaClientCapabilities
}

// Int32Response is a response to a Int32Request.
type Int32Response struct {
	// Diagnostics report errors or warnings related to validating the data source, provider, or resource
	// configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
