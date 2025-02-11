// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// CreateRequest represents a request for the provider to create a
// resource. An instance of this request struct is supplied as an argument to
// the resource's Create function.
type CreateRequest struct {
	// Config is the configuration the user supplied for the resource.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// Plan is the planned state for the resource.
	Plan tfsdk.Plan

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config
}

// CreateResponse represents a response to a CreateRequest. An
// instance of this response struct is supplied as
// an argument to the resource's Create function, in which the provider
// should set values on the CreateResponse as appropriate.
type CreateResponse struct {
	// State is the state of the resource following the Create operation.
	// This field is pre-populated from CreateRequest.Plan and
	// should be set during the resource's Create operation.
	State tfsdk.State

	// Private is the private state resource data following the Create operation.
	// This field is not pre-populated as there is no pre-existing private state
	// data during the resource's Create operation.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to creating the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
