// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"errors"
	"fmt"
)

var _ PlanCheck = expectEmptyPlan{}

type expectEmptyPlan struct{}

// CheckPlan implements the plan check logic.
func (e expectEmptyPlan) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var result []error

	for output, change := range req.Plan.OutputChanges {
		if !change.Actions.NoOp() {
			result = append(result, fmt.Errorf("expected empty plan, but output %q has planned action(s): %v", output, change.Actions))
		}
	}

	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			result = append(result, fmt.Errorf("expected empty plan, but %s has planned action(s): %v", rc.Address, rc.Change.Actions))
		}
	}

	resp.Error = errors.Join(result...)
}

// ExpectEmptyPlan returns a plan check that asserts that there are no output or resource changes in the plan.
// All output and resource changes found will be aggregated and returned in a plan check error.
func ExpectEmptyPlan() PlanCheck {
	return expectEmptyPlan{}
}
