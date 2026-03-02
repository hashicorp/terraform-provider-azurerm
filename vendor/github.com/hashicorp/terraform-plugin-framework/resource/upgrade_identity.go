// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Request information for the provider logic to update a resource identity
// from a prior resource identity version to the current identity version.
type UpgradeIdentityRequest struct {
	// Previous state of the resource identity in JSON format
	// (Terraform CLI 0.12 and later) This data is always available,
	// regardless of whether the wrapping IdentityUpgrader type
	// PriorSchema field was present.
	//
	// This is advanced functionality for providers wanting to skip the full
	// redeclaration of older identity schemas and instead use lower level handlers
	// to transform data. A typical implementation for working with this data will
	// call the Unmarshal() method.
	RawIdentity *tfprotov6.RawState

	// Previous identity of the resource if the wrapping IdentityUpgrader
	// type PriorSchema field was present. When available, this allows for
	// easier data handling such as calling Get() or GetAttribute().
	Identity *tfsdk.ResourceIdentity
}

// Response information for the provider logic to update a resource identity
// from a prior resource identity version to the current identity version.
type UpgradeIdentityResponse struct {
	// Upgraded identity of the resource, which should match the current identity
	//schema version.
	//
	// This field allows for easier data handling such as calling Set() or
	// SetAttribute().
	//
	// All data must be populated to prevent data loss during the upgrade
	// operation. No prior identity data is copied automatically.
	Identity *tfsdk.ResourceIdentity

	// Diagnostics report errors or warnings related to upgrading the resource
	// identity state. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
