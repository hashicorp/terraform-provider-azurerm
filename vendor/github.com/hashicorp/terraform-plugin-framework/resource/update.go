// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// UpdateRequest represents a request for the provider to update a
// resource. An instance of this request struct is supplied as an argument to
// the resource's Update function.
type UpdateRequest struct {
	// Config is the configuration the user supplied for the resource.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// Plan is the planned state for the resource.
	Plan tfsdk.Plan

	// State is the current state of the resource prior to the Update
	// operation.
	State tfsdk.State

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state. Any existing data is copied to
	// UpdateResponse.Private to prevent accidental private state data loss.
	//
	// Use the GetKey method to read data. Use the SetKey method on
	// UpdateResponse.Private to update or remove a value.
	Private *privatestate.ProviderData
}

// UpdateResponse represents a response to an UpdateRequest. An
// instance of this response struct is supplied as
// an argument to the resource's Update function, in which the provider
// should set values on the UpdateResponse as appropriate.
type UpdateResponse struct {
	// State is the state of the resource following the Update operation.
	// This field is pre-populated from UpdateResourceRequest.Plan and
	// should be set during the resource's Update operation.
	State tfsdk.State

	// Private is the private state resource data following the Update operation.
	// This field is pre-populated from UpdateRequest.Private and
	// can be modified during the resource's Update operation.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to updating the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
