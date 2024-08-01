// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ConfigureRequest represents a request containing the values the user
// specified for the provider configuration block, along with other runtime
// information from Terraform or the Plugin SDK. An instance of this request
// struct is supplied as an argument to the provider's Configure function.
type ConfigureRequest struct {
	// TerraformVersion is the version of Terraform executing the request.
	// This is supplied for logging, analytics, and User-Agent purposes
	// only. Providers should not try to gate provider behavior on
	// Terraform versions.
	TerraformVersion string

	// Config is the configuration the user supplied for the provider. This
	// information should usually be persisted to the underlying type
	// that's implementing the Provider interface, for use in later
	// resource CRUD operations.
	Config tfsdk.Config
}

// ConfigureResponse represents a response to a
// ConfigureRequest. An instance of this response struct is supplied as
// an argument to the provider's Configure function, in which the provider
// should set values on the ConfigureResponse as appropriate.
type ConfigureResponse struct {
	// DataSourceData is provider-defined data, clients, etc. that is passed
	// to [datasource.ConfigureRequest.ProviderData] for each DataSource type
	// that implements the Configure method.
	DataSourceData any

	// Diagnostics report errors or warnings related to configuring the
	// provider. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics

	// ResourceData is provider-defined data, clients, etc. that is passed
	// to [resource.ConfigureRequest.ProviderData] for each Resource type
	// that implements the Configure method.
	ResourceData any
}
