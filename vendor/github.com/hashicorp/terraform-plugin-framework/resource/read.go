// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ReadRequest represents a request for the provider to read a
// resource, i.e., update values in state according to the real state of the
// resource. An instance of this request struct is supplied as an argument to
// the resource's Read function.
type ReadRequest struct {
	// State is the current state of the resource prior to the Read
	// operation.
	State tfsdk.State

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state. This data is opaque to Terraform and does
	// not affect plan output. Any existing data is copied to
	// ReadResourceResponse.Private to prevent accidental private state data loss.
	//
	// Use the GetKey method to read data. Use the SetKey method on
	// ReadResourceResponse.Private to update or remove a value.
	Private *privatestate.ProviderData

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config
}

// ReadResponse represents a response to a ReadRequest. An
// instance of this response struct is supplied as
// an argument to the resource's Read function, in which the provider
// should set values on the ReadResponse as appropriate.
type ReadResponse struct {
	// State is the state of the resource following the Read operation.
	// This field is pre-populated from ReadRequest.State and
	// should be set during the resource's Read operation.
	State tfsdk.State

	// Private is the private state resource data following the Read operation.
	// This field is pre-populated from ReadResourceRequest.Private and
	// can be modified during the resource's Read operation.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to reading the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
