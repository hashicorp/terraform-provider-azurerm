// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ConfigureProviderClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ConfigureProvider RPC,
// such as forward-compatible Terraform behavior changes.
type ConfigureProviderClientCapabilities struct {
	// DeferralAllowed indicates whether the Terraform client initiating
	// the request allows a deferral response.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	DeferralAllowed bool
}

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

	// ClientCapabilities defines optionally supported protocol features for the
	// ConfigureProvider RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities ConfigureProviderClientCapabilities
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

	// EphemeralResourceData is provider-defined data, clients, etc. that is
	// passed to [ephemeral.ConfigureRequest.ProviderData] for each
	// EphemeralResource type that implements the Configure method.
	EphemeralResourceData any

	// ActionData is provider-defined data, clients, etc. that is
	// passed to [action.ConfigureRequest.ProviderData] for each
	// Action type that implements the Configure method.
	ActionData any

	// ListResourceData is provider-defined data, clients, etc. that is
	// passed to [action.ConfigureRequest.ProviderData] for each
	// Action type that implements the Configure method.
	ListResourceData any

	// Deferred indicates that Terraform should automatically defer
	// all resources and data sources for this provider.
	//
	// This field can only be set if
	// `(provider.ConfigureRequest).ClientCapabilities.DeferralAllowed` is true.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	Deferred *Deferred
}
