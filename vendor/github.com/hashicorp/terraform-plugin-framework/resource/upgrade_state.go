// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Request information for the provider logic to update a resource state
// from a prior state version to the current schema version. An instance of
// this is supplied as a parameter to a StateUpgrader, which ultimately comes
// from a Resource's UpgradeState method.
type UpgradeStateRequest struct {
	// Previous state of the resource in JSON (Terraform CLI 0.12 and later)
	// or flatmap format, depending on which version of Terraform CLI last
	// wrote the resource state. This data is always available, regardless
	// whether the wrapping StateUpgrader type PriorSchema field was
	// present.
	//
	// This is advanced functionality for providers wanting to skip the full
	// redeclaration of older schemas and instead use lower level handlers to
	// transform data. A typical implementation for working with this data will
	// call the Unmarshal() method.
	//
	// TODO: Create framework defined type that is not protocol specific.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/340
	RawState *tfprotov6.RawState

	// Previous state of the resource if the wrapping StateUpgrader
	// type PriorSchema field was present. When available, this allows for
	// easier data handling such as calling Get() or GetAttribute().
	State *tfsdk.State
}

// Response information for the provider logic to update a resource state
// from a prior state version to the current schema version. An instance of
// this is supplied as a parameter to a StateUpgrader, which ultimately came
// from a Resource's UpgradeState method.
type UpgradeStateResponse struct {
	// Diagnostics report errors or warnings related to upgrading the resource
	// state. An empty slice indicates a successful operation with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics

	// Upgraded state of the resource, which should match the current schema
	// version. If set, this will override State.
	//
	// This field is intended only for advanced provider functionality, such as
	// skipping the full redeclaration of older schemas or using lower level
	// handlers to transform data. Call tfprotov6.NewDynamicValue() to set this
	// value.
	//
	// All data must be populated to prevent data loss during the upgrade
	// operation. No prior state data is copied automatically.
	//
	// TODO: Remove in preference of requiring State, rather than using either
	// a new framework defined type or keeping this protocol specific type.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/340
	DynamicValue *tfprotov6.DynamicValue

	// Upgraded state of the resource, which should match the current schema
	// version. If DynamicValue is set, it will override this value.
	//
	// This field allows for easier data handling such as calling Set() or
	// SetAttribute(). It is generally recommended over working with the lower
	// level types and functionality required for DynamicValue.
	//
	// All data must be populated to prevent data loss during the upgrade
	// operation. No prior state data is copied automatically.
	State tfsdk.State
}
