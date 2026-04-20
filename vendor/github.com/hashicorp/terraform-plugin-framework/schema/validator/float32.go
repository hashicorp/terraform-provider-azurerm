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

// Float32 is a schema validator for types.Float32 attributes.
type Float32 interface {
	Describer

	// ValidateFloat32 should perform the validation.
	ValidateFloat32(context.Context, Float32Request, *Float32Response)
}

// Float32Request is a request for types.Float32 schema validation.
type Float32Request struct {
	// Path contains the path of the attribute for validation. Use this path
	// for any response diagnostics.
	Path path.Path

	// PathExpression contains the expression matching the exact path
	// of the attribute for validation.
	PathExpression path.Expression

	// Config contains the entire configuration of the data source, provider, or resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for validation from the configuration.
	ConfigValue types.Float32

	// ClientCapabilities defines optionally supported protocol features for
	// schema validation RPCs, such as forward-compatible Terraform
	// behavior changes.
	ClientCapabilities ValidateSchemaClientCapabilities
}

// Float32Response is a response to a Float32Request.
type Float32Response struct {
	// Diagnostics report errors or warnings related to validating the data source, provider, or resource
	// configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
