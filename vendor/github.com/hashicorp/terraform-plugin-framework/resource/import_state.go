// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportStateClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the ImportResourceState RPC,
// such as forward-compatible Terraform behavior changes.
type ImportStateClientCapabilities struct {
	// DeferralAllowed indicates whether the Terraform client initiating
	// the request allows a deferral response.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	DeferralAllowed bool
}

// ImportStateRequest represents a request for the provider to import a
// resource. An instance of this request struct is supplied as an argument to
// the Resource's ImportState method.
type ImportStateRequest struct {
	// ID represents the import identifier supplied by the practitioner when
	// calling the import command. In many cases, this may align with the
	// unique identifier for the resource, which can optionally be stored
	// as an Attribute. However, this identifier can also be treated as
	// its own type of value and parsed during import. This value
	// is not stored in the state unless the provider explicitly stores it.
	//
	// This ID field is supplied in the "terraform import" CLI command or in the import config block "id" attribute.
	// Either ID or Identity must be supplied by the practitioner, depending on the method used to import the resource.
	ID string

	// Identity is the configuration data provided by the practitioner in the import config block "identity" attribute. This
	// configuration data will conform to the identity schema defined by the managed resource. If the resource does not support identity,
	// this value will not be set.
	//
	// The "identity" attribute in the import block is only supported in Terraform 1.12 and later.
	// Either ID or Identity must be supplied by the practitioner, depending on the method used to import the resource.
	Identity *tfsdk.ResourceIdentity

	// ClientCapabilities defines optionally supported protocol features for the
	// ImportResourceState RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities ImportStateClientCapabilities
}

// ImportStateResponse represents a response to a ImportStateRequest.
// An instance of this response struct is supplied as an argument to the
// Resource's ImportState method, in which the provider should set values on
// the ImportStateResponse as appropriate.
type ImportStateResponse struct {
	// Diagnostics report errors or warnings related to importing the
	// resource. An empty slice indicates a successful operation with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics

	// State is the state of the resource following the import operation.
	// It must contain enough information so Terraform can successfully
	// refresh the resource, e.g. call the Resource Read method.
	State tfsdk.State

	// Identity is the identity of the resource following the Import operation.
	// This field is pre-populated from ImportStateRequest.Identity and
	// should be set during the resource's Import operation.
	//
	// If the resource does not support identity, this value will not be set and will
	// raise a diagnostic if set by the resource's Import operation.
	Identity *tfsdk.ResourceIdentity

	// Private is the private state resource data following the Import operation.
	// This field is not pre-populated as there is no pre-existing private state
	// data during the resource's Import operation.
	Private *privatestate.ProviderData

	// Deferred indicates that Terraform should defer
	// importing this resource.
	//
	// This field can only be set if
	// `(resource.ImportStateRequest).ClientCapabilities.DeferralAllowed` is true.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	Deferred *Deferred
}

// ImportStatePassthroughID is a helper function to set the import
// identifier to a given state attribute path. The attribute must accept a
// string value.
//
// For resources that support identity, this method will also automatically pass through the
// Identity field if imported by the identity attribute of a import config block (Terraform 1.12+ and later).
// In this scenario where identity is provided instead of the string ID, the state field defined
// at `attrPath` will be set to null.
func ImportStatePassthroughID(ctx context.Context, attrPath path.Path, req ImportStateRequest, resp *ImportStateResponse) {
	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughID path must be set to a valid attribute path that can accept a string value.",
		)
		return
	}

	// If the import is using the ID string identifier, (either via the "terraform import" CLI command, or a config block with the "id" attribute set)
	// pass through the ID to the designated state attribute.
	if req.ID != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, req.ID)...)
	}
}

// ImportStatePassthroughWithIdentity is a helper function to retrieve either the import identifier
// or a given identity attribute that is then used to set to given attribute path in state, based on the method used
// by the practitioner to import. The identity and state attributes provided must be of type string.
//
// The helper method should only be used on resources that support identity via the resource.ResourceWithIdentity interface.
//
// This method will also automatically pass through the Identity field if imported by
// the identity attribute of a import config block (Terraform 1.12+ and later).
func ImportStatePassthroughWithIdentity(ctx context.Context, stateAttrPath, identityAttrPath path.Path, req ImportStateRequest, resp *ImportStateResponse) {
	if stateAttrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing State Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughWithIdentity path must be set to a valid state attribute path that can accept a string value.",
		)
	}

	if identityAttrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Identity Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughWithIdentity path must be set to a valid identity attribute path that is a string value.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If the import is using the import identifier, (either via the "terraform import" CLI command, or a config block with the "id" attribute set)
	// pass through the ID to the designated state attribute.
	if req.ID != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, stateAttrPath, req.ID)...)
		return
	}

	// The import isn't using the import identifier, so it must be using identity. Grab the designated
	// identity attribute string and set it to state.
	var identityAttrVal types.String
	resp.Diagnostics.Append(req.Identity.GetAttribute(ctx, identityAttrPath, &identityAttrVal)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, stateAttrPath, identityAttrVal)...)
}
