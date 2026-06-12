// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
)

// ListResourceSchemaRequest represents a request for the ListResource to
// return its schema.  An instance of this request struct is supplied as an
// argument to the ListResource type ListResourceSchema method.
type ListResourceSchemaRequest struct{}

// ListResourceSchemaResponse represents a response to a
// ListResourceSchemaRequest. An instance of this response struct is supplied
// as an argument to the ListResource type ListResourceResourceSchema method.
type ListResourceSchemaResponse struct {
	// Schema is the schema of the list resource.
	Schema schema.Schema

	// Diagnostics report errors or warnings related to retrieving the list
	// resource schema. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
