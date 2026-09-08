// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// DeleteStateRequest represents a request to delete the given state ([DeleteStateRequest.StateID]) in the configured state store.
type DeleteStateRequest struct {
	// StateID is the ID of the state to delete.
	//
	// Typically this is the name of the Terraform workspace the practitioner is
	// running Terraform in: https://developer.hashicorp.com/terraform/language/state/workspaces .
	//
	// If the practitioner hasn't explicitly selected a workspace, StateID will be set to "default".
	StateID string
}

// DeleteStateResponse represents a response to an DeleteStateRequest. An instance of this response
// struct is supplied as an argument to the state store's DeleteState method, in which the provider
// should set values on the DeleteStateResponse as appropriate.
type DeleteStateResponse struct {
	// Diagnostics report errors or warnings related to deleting a state in the configured
	// state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
