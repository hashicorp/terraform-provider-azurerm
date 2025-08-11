// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
)

// IdentitySchemaRequest represents a request for the Resource to return its identity schema.
// An instance of this request struct is supplied as an argument to the
// Resource type IdentitySchema method.
type IdentitySchemaRequest struct{}

// IdentitySchemaResponse represents a response to a SchemaRequest. An instance of this
// response struct is supplied as an argument to the Resource type IdentitySchema
// method.
type IdentitySchemaResponse struct {
	// IdentitySchema is the schema of the resource identity.
	IdentitySchema identityschema.Schema

	// Diagnostics report errors or warnings related to retrieving the resource
	// identity schema. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
