// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// WriteRequest represents a request to write the state data ([WriteRequest.StateBytes]) to
// a given state ([WriteRequest.StateID]) in the state store.
type WriteRequest struct {
	// StateID is the ID of the state to write to.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces .
	//
	// If the practitioner hasn't explicitly selected a workspace, StateID will be set to "default".
	StateID string

	// StateBytes is the entire state file to write to [WriteRequest.StateID].
	StateBytes []byte
}

// WriteResponse represents a response to a WriteRequest. An instance of this response
// struct is supplied as an argument to the state store's Write method, in which the provider
// should set values on the WriteResponse as appropriate.
type WriteResponse struct {
	// Diagnostics report errors or warnings related to writing state data to a
	// state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
