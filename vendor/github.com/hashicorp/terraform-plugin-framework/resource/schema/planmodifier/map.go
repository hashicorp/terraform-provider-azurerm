// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Map is a schema plan modifier for types.Map attributes.
type Map interface {
	Describer

	// PlanModifyMap should perform the modification.
	PlanModifyMap(context.Context, MapRequest, *MapResponse)
}

// MapRequest is a request for types.Map schema plan modification.
type MapRequest struct {
	// Path contains the path of the attribute for modification. Use this path
	// for any response diagnostics.
	Path path.Path

	// PathExpression contains the expression matching the exact path
	// of the attribute for modification.
	PathExpression path.Expression

	// Config contains the entire configuration of the resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for modification from the configuration.
	ConfigValue types.Map

	// Plan contains the entire proposed new state of the resource.
	Plan tfsdk.Plan

	// PlanValue contains the value of the attribute for modification from the proposed new state.
	PlanValue types.Map

	// State contains the entire prior state of the resource.
	State tfsdk.State

	// StateValue contains the value of the attribute for modification from the prior state.
	StateValue types.Map

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state. This data is opaque to Terraform and does
	// not affect plan output. Any existing data is copied to
	// MapResponse.Private to prevent accidental private state data loss.
	//
	// The private state data is always the original data when the schema-based plan
	// modification began or, is updated as the logic traverses deeper into underlying
	// attributes.
	//
	// Use the GetKey method to read data. Use the SetKey method on
	// MapResponse.Private to update or remove a value.
	Private *privatestate.ProviderData
}

// MapResponse is a response to a MapRequest.
type MapResponse struct {
	// PlanValue is the planned new state for the attribute.
	PlanValue types.Map

	// RequiresReplace indicates whether a change in the attribute
	// requires replacement of the whole resource.
	RequiresReplace bool

	// Private is the private state resource data following the PlanModifyMap operation.
	// This field is pre-populated from MapRequest.Private and
	// can be modified during the resource's PlanModifyMap operation.
	//
	// The private state data is always the original data when the schema-based plan
	// modification began or, is updated as the logic traverses deeper into underlying
	// attributes.
	Private *privatestate.ProviderData

	// Diagnostics report errors or warnings related to modifying the resource
	// plan. An empty slice indicates success, with no warnings or
	// errors generated.
	Diagnostics diag.Diagnostics
}
