// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectUnknownOutputValue{}

type expectUnknownOutputValue struct {
	outputAddress string
}

// CheckPlan implements the plan check logic.
func (e expectUnknownOutputValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
	var change *tfjson.Change

	for address, oc := range req.Plan.OutputChanges {
		if e.outputAddress == address {
			change = oc

			break
		}
	}

	if change == nil {
		resp.Error = fmt.Errorf("%s - Output not found in plan OutputChanges", e.outputAddress)

		return
	}

	result, err := tfjsonpath.Traverse(change.AfterUnknown, tfjsonpath.Path{})
	if err != nil {
		// If we find the output in the known values, return a more explicit message
		knownVal, knownErr := tfjsonpath.Traverse(change.After, tfjsonpath.Path{})
		if knownErr == nil {
			resp.Error = fmt.Errorf("Expected unknown value at output %q, but found known value: \"%v\"", e.outputAddress, knownVal)
			return
		}

		resp.Error = err
		return
	}

	isUnknown, ok := result.(bool)

	if !ok {
		resp.Error = fmt.Errorf("invalid path: the path value cannot be asserted as bool")

		return
	}

	if !isUnknown {
		// The output should have a known value, look first to return a more explicit message
		knownVal, knownErr := tfjsonpath.Traverse(change.After, tfjsonpath.Path{})
		if knownErr == nil {
			resp.Error = fmt.Errorf("Expected unknown value at output %q, but found known value: \"%v\"", e.outputAddress, knownVal)
			return
		}
		resp.Error = fmt.Errorf("Expected unknown value at output %q, but found known value", e.outputAddress)

		return
	}
}

// ExpectUnknownOutputValue returns a plan check that asserts that the specified output has an unknown value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of unknown
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of unknown values, such
// as marking whole maps as unknown rather than individual element values.
func ExpectUnknownOutputValue(outputAddress string) PlanCheck {
	return expectUnknownOutputValue{
		outputAddress: outputAddress,
	}
}
