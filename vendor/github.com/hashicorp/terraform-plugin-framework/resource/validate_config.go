// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateConfigClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the
// ValidateResourceConfig RPC, such as forward-compatible Terraform behavior
// changes.
type ValidateConfigClientCapabilities struct {
	// WriteOnlyAttributesAllowed indicates that the Terraform client
	// initiating the request supports write-only attributes for managed
	// resources.
	WriteOnlyAttributesAllowed bool
}

// ValidateConfigRequest represents a request to validate the
// configuration of a resource. An instance of this request struct is
// supplied as an argument to the Resource ValidateConfig receiver method
// or automatically passed through to each ConfigValidator.
type ValidateConfigRequest struct {
	// Config is the configuration the user supplied for the resource.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// ClientCapabilities defines optionally supported protocol features for
	// the ValidateResourceConfig RPC, such as forward-compatible Terraform
	// behavior changes.
	ClientCapabilities ValidateConfigClientCapabilities
}

// ValidateConfigResponse represents a response to a
// ValidateConfigRequest. An instance of this response struct is
// supplied as an argument to the Resource ValidateConfig receiver method
// or automatically passed through to each ConfigValidator.
type ValidateConfigResponse struct {
	// Diagnostics report errors or warnings related to validating the resource
	// configuration. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
