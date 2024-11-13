// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ReadClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ReadDataSource RPC,
// such as forward-compatible Terraform behavior changes.
type ReadClientCapabilities struct {
	// DeferralAllowed indicates whether the Terraform client initiating
	// the request allows a deferral response.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	DeferralAllowed bool
}

// ReadRequest represents a request for the provider to read a data
// source, i.e., update values in state according to the real state of the
// data source. An instance of this request struct is supplied as an argument
// to the data source's Read function.
type ReadRequest struct {
	// Config is the configuration the user supplied for the data source.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config

	// ClientCapabilities defines optionally supported protocol features for the
	// ReadDataSource RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities ReadClientCapabilities
}

// ReadResponse represents a response to a ReadRequest. An
// instance of this response struct is supplied as an argument to the data
// source's Read function, in which the provider should set values on the
// ReadResponse as appropriate.
type ReadResponse struct {
	// State is the state of the data source following the Read operation.
	// This field should be set during the resource's Read operation.
	State tfsdk.State

	// Diagnostics report errors or warnings related to reading the data
	// source. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics

	// Deferred indicates that Terraform should defer reading this
	// data source until a followup apply operation.
	//
	// This field can only be set if
	// `(datasource.ReadRequest).ClientCapabilities.DeferralAllowed` is true.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	Deferred *Deferred
}
