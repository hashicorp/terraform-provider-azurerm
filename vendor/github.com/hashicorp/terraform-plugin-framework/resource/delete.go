// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DeleteRequest represents a request for the provider to delete a
// resource. An instance of this request struct is supplied as an argument to
// the resource's Delete function.
type DeleteRequest struct {
	// State is the current state of the resource prior to the Delete
	// operation.
	State tfsdk.State

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state.
	//
	// Use the GetKey method to read data.
	Private *privatestate.ProviderData
}

// DeleteResponse represents a response to a DeleteRequest. An
// instance of this response struct is supplied as
// an argument to the resource's Delete function, in which the provider
// should set values on the DeleteResponse as appropriate.
type DeleteResponse struct {
	// State is the state of the resource following the Delete operation.
	// This field is pre-populated from UpdateResourceRequest.Plan and
	// should be set during the resource's Update operation.
	State tfsdk.State

	// Private is the private state resource data following the Delete
	// operation. This field is pre-populated from DeleteRequest.Private and
	// can be modified during the resource's Delete operation in cases where
	// an error diagnostic is being returned. Otherwise if no error diagnostic
	// is being returned, indicating that the resource was successfully deleted,
	// this data will be automatically cleared to prevent Terraform errors.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to deleting the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
