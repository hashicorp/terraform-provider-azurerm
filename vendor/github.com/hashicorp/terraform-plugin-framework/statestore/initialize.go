// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package statestore

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// InitializeRequest represents a request containing the values the user
// specified for the state_store configuration block, along with the data configured
// for the provider itself, populated by the [provider.ConfigureResponse.StateStoreData] field.
//
// An instance of this request struct is supplied as an argument to the state store's
// Initialize method.
type InitializeRequest struct {
	// Config is the configuration the user supplied for the state store.
	Config tfsdk.Config

	// ProviderData is the data set in the [provider.ConfigureResponse.StateStoreData]
	// field. This data is provider-specifc and therefore can contain any necessary remote system
	// clients, custom provider data, or anything else pertinent to the functionality of the StateStore.
	ProviderData any
}

// InitializeResponse represents a response to an InitializeRequest. An instance of this response
// struct is supplied as an argument to the state store's Initialize method, in which the provider
// should set values on the InitializeResponse as appropriate.
type InitializeResponse struct {
	// Diagnostics report errors or warnings related to initializing the
	// state store. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics

	// StateStoreData is a combination of provider-defined and state store specific data that is
	// passed to [ConfigureRequest.StateStoreData] for the implementing StateStore type.
	StateStoreData any
}
