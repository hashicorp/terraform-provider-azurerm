// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"errors"
)

var _ PlanCheck = expectNonEmptyPlan{}

type expectNonEmptyPlan struct{}

// CheckPlan implements the plan check logic.
func (e expectNonEmptyPlan) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			return
		}
	}

	resp.Error = errors.New("expected a non-empty plan, but got an empty plan")
}

// ExpectNonEmptyPlan returns a plan check that asserts there is at least one resource change in the plan.
func ExpectNonEmptyPlan() PlanCheck {
	return expectNonEmptyPlan{}
}
