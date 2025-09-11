// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package action

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ModifyPlanClientCapabilities allows Terraform to publish information
// regarding optionally supported protocol features for the PlanAction RPC,
// such as forward-compatible Terraform behavior changes.
type ModifyPlanClientCapabilities struct {
	// DeferralAllowed indicates whether the Terraform client initiating
	// the request allows a deferral response.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	DeferralAllowed bool
}

// ModifyPlanRequest represents a request for the provider to modify the
// planned new state that Terraform has generated for any linked resources.
type ModifyPlanRequest struct {
	// Config is the configuration the user supplied for the action.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// TODO:Actions: Add linked resources once lifecycle/linked actions are implemented

	// ClientCapabilities defines optionally supported protocol features for the
	// PlanAction RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities ModifyPlanClientCapabilities
}

// ModifyPlanResponse represents a response to a
// ModifyPlanRequest. An instance of this response struct is supplied
// as an argument to the action's ModifyPlan function, in which the provider
// should modify the Plan of any linked resources as appropriate.
type ModifyPlanResponse struct {
	// Diagnostics report errors or warnings related to determining the
	// planned state of the requested action's linked resources. Returning an empty slice
	// indicates a successful plan modification with no warnings or errors
	// generated.
	Diagnostics diag.Diagnostics

	// TODO:Actions: Add linked resources once lifecycle/linked actions are implemented

	// Deferred indicates that Terraform should defer planning this
	// action until a follow-up apply operation.
	//
	// This field can only be set if
	// `(action.ModifyPlanRequest).ClientCapabilities.DeferralAllowed` is true.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	Deferred *Deferred
}
