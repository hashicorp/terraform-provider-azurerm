// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// GetStatesRequest represents a request to retrieve all state IDs for states persisted in the configured state store.
type GetStatesRequest struct{}

// GetStatesResponse represents a response to a GetStatesRequest. An instance of this response
// struct is supplied as an argument to the state store's GetStates method, in which the provider
// should set values on the GetStatesResponse as appropriate.
type GetStatesResponse struct {
	// StateIDs is a list of all state IDs for states persisted in the configured state store.
	StateIDs []string

	// Diagnostics report errors or warnings related to retrieving all state IDs for states
	// persisted in the configured state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
