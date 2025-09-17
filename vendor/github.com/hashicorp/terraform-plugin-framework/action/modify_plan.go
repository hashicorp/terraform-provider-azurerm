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

// ModifyPlanRequest represents a request for the provider during planning.
// The plan can be used as an opportunity to raise early
// diagnostics to practitioners, such as validation errors.
type ModifyPlanRequest struct {
	// Config is the configuration the user supplied for the action.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config

	// ClientCapabilities defines optionally supported protocol features for the
	// PlanAction RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities ModifyPlanClientCapabilities
}

// ModifyPlanResponse represents a response to a
// ModifyPlanRequest. An instance of this response struct is supplied
// as an argument to the action's ModifyPlan function.
type ModifyPlanResponse struct {
	// Diagnostics report early errors or warnings related action.
	// Returning an empty slice indicates a successful plan modification
	// with no warnings or errors generated.
	Diagnostics diag.Diagnostics

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
