// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseNonNullStateForUnknown returns a plan modifier that copies a known, non-null, prior state
// value into the planned value. Use this when it is known that an unconfigured value will remain the
// same after the attribute is updated to a non-null value.
//
// To prevent Terraform errors, the framework automatically sets unconfigured
// and Computed attributes to an unknown value "(known after apply)" on update.
// Using this plan modifier will instead display the non-null prior state value in the
// plan, unless a prior plan modifier adjusts the value.
//
// This plan modifier can be a useful alternative to [UseStateForUnknown] when the attribute is
// a child of a nested attribute that can be null after the resource is created.
func UseNonNullStateForUnknown() planmodifier.Set {
	return useNonNullStateForUnknown{}
}

type useNonNullStateForUnknown struct{}

func (m useNonNullStateForUnknown) Description(_ context.Context) string {
	return "Once set to a non-null value, the value of this attribute in state will not change."
}

func (m useNonNullStateForUnknown) MarkdownDescription(_ context.Context) string {
	return "Once set to a non-null value, the value of this attribute in state will not change."
}

func (m useNonNullStateForUnknown) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Do nothing if the state value is null.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
