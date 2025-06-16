// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
)

// CloseRequest represents a request for the provider to close an ephemeral
// resource. An instance of this request struct is supplied as an argument to
// the ephemeral resource's Close function.
type CloseRequest struct {
	// Private is provider-defined ephemeral resource private state data
	// which was previously provided by the latest Open or Renew operation.
	//
	// Use the GetKey method to read data.
	Private *privatestate.ProviderData
}

// CloseResponse represents a response to a CloseRequest. An
// instance of this response struct is supplied as an argument
// to the ephemeral resource's Close function, in which the provider
// should set values on the CloseResponse as appropriate.
type CloseResponse struct {
	// Diagnostics report errors or warnings related to closing the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
