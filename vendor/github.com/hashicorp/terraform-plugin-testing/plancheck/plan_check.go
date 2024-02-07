// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// PlanCheck defines an interface for implementing test logic that checks a plan file and then returns an error
// if the plan file does not match what is expected.
type PlanCheck interface {
	// CheckPlan should perform the plan check.
	CheckPlan(context.Context, CheckPlanRequest, *CheckPlanResponse)
}

// CheckPlanRequest is a request for an invoke of the CheckPlan function.
type CheckPlanRequest struct {
	// Plan represents a parsed plan file, retrieved via the `terraform show -json` command.
	Plan *tfjson.Plan
}

// CheckPlanResponse is a response to an invoke of the CheckPlan function.
type CheckPlanResponse struct {
	// Error is used to report the failure of a plan check assertion and is combined with other PlanCheck errors
	// to be reported as a test failure.
	Error error
}
