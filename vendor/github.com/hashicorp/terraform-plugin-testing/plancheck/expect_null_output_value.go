// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plancheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ PlanCheck = expectNullOutputValue{}

type expectNullOutputValue struct {
	outputAddress string
}

// CheckPlan implements the plan check logic.
func (e expectNullOutputValue) CheckPlan(ctx context.Context, req CheckPlanRequest, resp *CheckPlanResponse) {
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

	var result any
	var err error

	switch {
	case change.Actions.Create():
		result, err = tfjsonpath.Traverse(change.After, tfjsonpath.Path{})
	default:
		result, err = tfjsonpath.Traverse(change.Before, tfjsonpath.Path{})
	}

	if err != nil {
		resp.Error = err

		return
	}

	if result != nil {
		resp.Error = fmt.Errorf("attribute at path is not null")

		return
	}
}

// ExpectNullOutputValue returns a plan check that asserts that the specified output has a null value.
//
// Due to implementation differences between the terraform-plugin-sdk and the terraform-plugin-framework, representation of null
// values may differ. For example, terraform-plugin-sdk based providers may have less precise representations of null values, such
// as marking whole maps as null rather than individual element values.
//
// Deprecated: Use [plancheck.ExpectKnownOutputValue] with [knownvalue.Null] instead.
// ExpectNullOutputValue will be removed in the next major version release.
func ExpectNullOutputValue(outputAddress string) PlanCheck {
	return expectNullOutputValue{
		outputAddress: outputAddress,
	}
}
