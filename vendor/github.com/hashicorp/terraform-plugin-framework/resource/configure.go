// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ConfigureRequest represents a request for the provider to configure a
// resource, i.e., set provider-level data or clients. An instance of this
// request struct is supplied as an argument to the Resource type Configure
// method.
type ConfigureRequest struct {
	// ProviderData is the data set in the
	// [provider.ConfigureResponse.ResourceData] field. This data is
	// provider-specifc and therefore can contain any necessary remote system
	// clients, custom provider data, or anything else pertinent to the
	// functionality of the Resource.
	//
	// This data is only set after the ConfigureProvider RPC has been called
	// by Terraform.
	ProviderData any
}

// ConfigureResponse represents a response to a ConfigureRequest. An
// instance of this response struct is supplied as an argument to the
// Resource type Configure method.
type ConfigureResponse struct {
	// Diagnostics report errors or warnings related to configuring of the
	// Datasource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
