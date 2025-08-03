// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// OpenClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the OpenEphemeralResource RPC,
// such as forward-compatible Terraform behavior changes.
type OpenClientCapabilities struct {
	// DeferralAllowed indicates whether the Terraform client initiating
	// the request allows a deferral response.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	DeferralAllowed bool
}

// OpenRequest represents a request for the provider to open an ephemeral
// resource. An instance of this request struct is supplied as an argument to
// the ephemeral resource's Open function.
type OpenRequest struct {
	// Config is the configuration the user supplied for the ephemeral
	// resource.
	Config tfsdk.Config

	// ClientCapabilities defines optionally supported protocol features for the
	// OpenEphemeralResource RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities OpenClientCapabilities
}

// OpenResponse represents a response to a OpenRequest. An
// instance of this response struct is supplied as an argument
// to the ephemeral resource's Open function, in which the provider
// should set values on the OpenResponse as appropriate.
type OpenResponse struct {
	// Result is the object representing the values of the ephemeral
	// resource following the Open operation. This field is pre-populated
	// from OpenRequest.Config and should be set during the resource's Open
	// operation.
	Result tfsdk.EphemeralResultData

	// Private is the private state ephemeral resource data following the
	// Open operation. This field is not pre-populated as there is no
	// pre-existing private state data during the ephemeral resource's
	// Open operation.
	//
	// This private data will be passed to any Renew or Close operations.
	Private *privatestate.ProviderData

	// RenewAt is an optional date/time field that indicates to Terraform
	// when this ephemeral resource must be renewed at. Terraform will call
	// the (EphemeralResource).Renew method when the current date/time is on
	// or after RenewAt during a Terraform operation.
	//
	// It is recommended to add extra time (usually no more than a few minutes)
	// before an ephemeral resource expires to account for latency.
	RenewAt time.Time

	// Diagnostics report errors or warnings related to opening the ephemeral
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics

	// Deferred indicates that Terraform should defer opening this
	// ephemeral resource until a followup apply operation.
	//
	// This field can only be set if
	// `(ephemeral.OpenRequest).ClientCapabilities.DeferralAllowed` is true.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	Deferred *Deferred
}
