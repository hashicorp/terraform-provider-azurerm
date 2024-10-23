// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

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
	ID string
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

	// Private is the private state resource data following the Import operation.
	// This field is not pre-populated as there is no pre-existing private state
	// data during the resource's Import operation.
	Private *privatestate.ProviderData
}

// ImportStatePassthroughID is a helper function to set the import
// identifier to a given state attribute path. The attribute must accept a
// string value.
func ImportStatePassthroughID(ctx context.Context, attrPath path.Path, req ImportStateRequest, resp *ImportStateResponse) {
	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughID path must be set to a valid attribute path that can accept a string value.",
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, req.ID)...)
}
