// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"
)

var _ PlanCheck = expectDeferredChange{}

type expectDeferredChange struct {
	resourceAddress string
	reason          DeferredReason
}

// CheckPlan implements the plan check logic.
func (e expectDeferredChange) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	foundResource := false

	for _, dc := range req.Plan.DeferredChanges {
		if dc.ResourceChange == nil || e.resourceAddress != dc.ResourceChange.Address {
			continue
		}

		if e.reason != DeferredReason(dc.Reason) {
			resp.Error = fmt.Errorf("'%s' - expected %q, got deferred reason: %q", dc.ResourceChange.Address, e.reason, dc.Reason)
			return
		}

		foundResource = true
		break
	}

	if !foundResource {
		resp.Error = fmt.Errorf("%s - No deferred changes found for resource", e.resourceAddress)
		return
	}
}

// ExpectDeferredChange returns a plan check that asserts that a given resource will have a
// deferred change in the plan with the given reason.
func ExpectDeferredChange(resourceAddress string, reason DeferredReason) PlanCheck {
	return expectDeferredChange{
		resourceAddress: resourceAddress,
		reason:          reason,
	}
}
