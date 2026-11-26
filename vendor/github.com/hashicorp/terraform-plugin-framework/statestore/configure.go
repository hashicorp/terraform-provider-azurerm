// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ConfigureRequest represents a request for the provider to set provider-level data or clients on a
// StateStore type. An instance of this request struct is supplied as an argument to the StateStore type
// Configure method.
//
// NOTE: This method is called any time framework needs to execute logic on a StateStore type and is different from
// the ConfigureStateStore RPC call, which is implemented by [StateStore.Initialize].
type ConfigureRequest struct {
	// StateStoreData is the data set in the [InitializeResponse.StateStoreData] field.
	// This data can contain any necessary remote system clients, state store initialization data,
	// or anything else pertinent to the functionality of the StateStore.
	//
	// This data is only set after the ConfigureStateStore RPC has been called by Terraform.
	StateStoreData any
}

// ConfigureResponse represents a response to a ConfigureRequest. An
// instance of this response struct is supplied as an argument to the
// StateStore type Configure method.
type ConfigureResponse struct {
	// Diagnostics report errors or warnings related to configuring of the
	// StateStore. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
