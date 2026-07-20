// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ReadRequest represents a request to read the contents of a given state ([ReadRequest.StateID]) in the state store.
type ReadRequest struct {
	// StateID is the ID of the state to read.
	//
	// Typically, this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces .
	//
	// If the practitioner hasn't explicitly selected a workspace, StateID will be set to "default".
	StateID string
}

// ReadResponse represents a response to an ReadRequest. An instance of this response
// struct is supplied as an argument to the state store's Read method, in which the provider
// should set values on the ReadResponse as appropriate.
type ReadResponse struct {
	// Diagnostics report errors or warnings related to reading the given state
	// from the state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics

	// StateBytes is the entire state file for [ReadRequest.StateID].
	StateBytes []byte
}
