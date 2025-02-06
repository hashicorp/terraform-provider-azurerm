// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource Plan Check
var _ PlanCheck = expectKnownValue{}

type expectKnownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	knownValue      knownvalue.Check
}

// CheckPlan implements the plan check logic.
func (e expectKnownValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var rc *tfjson.ResourceChange

	if req.Plan == nil {
		resp.Error = fmt.Errorf("plan is nil")
	}

	for _, resourceChange := range req.Plan.ResourceChanges {
		if e.resourceAddress == resourceChange.Address {
			rc = resourceChange

			break
		}
	}

	if rc == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in plan", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(rc.Change.After, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)

		return
	}
}

// ExpectKnownValue returns a plan check that asserts that the specified attribute at the given resource
// has a known type and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) PlanCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
