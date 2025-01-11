// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = &compareValuePairs{}

type compareValuePairs struct {
	resourceAddressOne string
	attributePathOne   tfjsonpath.Path
	resourceAddressTwo string
	attributePathTwo   tfjsonpath.Path
	comparer           compare.ValueComparer
}

// CheckState implements the state check logic.
func (e *compareValuePairs) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var resourceOne *tfjson.StateResource
	var resourceTwo *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddressOne == r.Address {
			resourceOne = r

			break
		}
	}

	if resourceOne == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressOne)

		return
	}

	resultOne, err := tfjsonpath.Traverse(resourceOne.AttributeValues, e.attributePathOne)

	if err != nil {
		resp.Error = err

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddressTwo == r.Address {
			resourceTwo = r

			break
		}
	}

	if resourceTwo == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddressTwo)

		return
	}

	resultTwo, err := tfjsonpath.Traverse(resourceTwo.AttributeValues, e.attributePathTwo)

	if err != nil {
		resp.Error = err

		return
	}

	err = e.comparer.CompareValues(resultOne, resultTwo)

	if err != nil {
		resp.Error = err
	}
}

// CompareValuePairs returns a state check that compares the value in state for the first given resource address and
// path with the value in state for the second given resource address and path using the supplied value comparer.
func CompareValuePairs(resourceAddressOne string, attributePathOne tfjsonpath.Path, resourceAddressTwo string, attributePathTwo tfjsonpath.Path, comparer compare.ValueComparer) StateCheck {
	return &compareValuePairs{
		resourceAddressOne: resourceAddressOne,
		attributePathOne:   attributePathOne,
		resourceAddressTwo: resourceAddressTwo,
		attributePathTwo:   attributePathTwo,
		comparer:           comparer,
	}
}
