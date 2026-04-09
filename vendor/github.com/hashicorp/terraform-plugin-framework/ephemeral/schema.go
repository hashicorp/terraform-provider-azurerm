// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
)

// SchemaRequest represents a request for the EphemeralResource to return its schema.
// An instance of this request struct is supplied as an argument to the
// EphemeralResource type Schema method.
type SchemaRequest struct{}

// SchemaResponse represents a response to a SchemaRequest. An instance of this
// response struct is supplied as an argument to the EphemeralResource type Schema
// method.
type SchemaResponse struct {
	// Schema is the schema of the ephemeral resource.
	Schema schema.Schema

	// Diagnostics report errors or warnings related to retrieving the ephemeral
	// resource schema. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
