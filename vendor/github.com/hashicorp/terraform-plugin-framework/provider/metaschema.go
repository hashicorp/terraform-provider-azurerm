// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
)

// MetaSchemaRequest represents a request for the Provider to return its schema.
// An instance of this request struct is supplied as an argument to the
// Provider type Schema method.
type MetaSchemaRequest struct{}

// MetaSchemaResponse represents a response to a MetaSchemaRequest. An instance of this
// response struct is supplied as an argument to the Provider type Schema
// method.
type MetaSchemaResponse struct {
	// Schema is the meta schema of the provider.
	Schema metaschema.Schema

	// Diagnostics report errors or warnings related to validating the data
	// source configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
