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
var _ PlanCheck = expectKnownOutputValueAtPath{}

type expectKnownOutputValueAtPath struct {
	outputAddress string
	outputPath    tfjsonpath.Path
	knownValue    knownvalue.Check
}

// CheckPlan implements the plan check logic.
func (e expectKnownOutputValueAtPath) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var change *tfjson.Change

	if req.Plan == nil {
		resp.Error = fmt.Errorf("plan is nil")
	}

	for address, oc := range req.Plan.OutputChanges {
		if e.outputAddress == address {
			change = oc

			break
		}
	}

	if change == nil {
		resp.Error = fmt.Errorf("%s - Output not found in plan", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(change.After, e.outputPath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for output at path: %s.%s, err: %s", e.outputAddress, e.outputPath.String(), err)

		return
	}
}

// ExpectKnownOutputValueAtPath returns a plan check that asserts that the specified output at the given path
// has a known type and value. Prior to Terraform v1.3.0 a planned output is marked as fully unknown
// if any attribute is unknown.
func ExpectKnownOutputValueAtPath(outputAddress string, outputPath tfjsonpath.Path, knownValue knownvalue.Check) PlanCheck {
	return expectKnownOutputValueAtPath{
		outputAddress: outputAddress,
		outputPath:    outputPath,
		knownValue:    knownValue,
	}
}
