// Copyright IBM Corp. 2014, 2025
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
var _ StateCheck = &compareValue{}

type compareValue struct {
	resourceAddresses []string
	attributePaths    []tfjsonpath.Path
	stateValues       []any
	comparer          compare.ValueComparer
}

func (e *compareValue) AddStateValue(resourceAddress string, attributePath tfjsonpath.Path) StateCheck {
	e.resourceAddresses = append(e.resourceAddresses, resourceAddress)
	e.attributePaths = append(e.attributePaths, attributePath)

	return e
}

// CheckState implements the state check logic.
func (e *compareValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
	var resource *tfjson.StateResource

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

	// All calls to AddStateValue occur before any TestStep is run, populating the resourceAddresses
	// and attributePaths slices. The stateValues slice is populated during execution of each TestStep.
	// Each call to CheckState happens sequentially during each TestStep.
	// The currentIndex is reflective of the current state value being checked.
	currentIndex := len(e.stateValues)

	if len(e.resourceAddresses) <= currentIndex {
		resp.Error = fmt.Errorf("resource addresses index out of bounds: %d", currentIndex)

		return
	}

	resourceAddress := e.resourceAddresses[currentIndex]

	for _, r := range req.State.Values.RootModule.Resources {
		if resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", resourceAddress)

		return
	}

	if len(e.attributePaths) <= currentIndex {
		resp.Error = fmt.Errorf("attribute paths index out of bounds: %d", currentIndex)

		return
	}

	attributePath := e.attributePaths[currentIndex]

	result, err := tfjsonpath.Traverse(resource.AttributeValues, attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	e.stateValues = append(e.stateValues, result)

	err = e.comparer.CompareValues(e.stateValues...)

	if err != nil {
		resp.Error = err
	}
}

// CompareValue returns a state check that compares values retrieved from state using the
// supplied value comparer.
func CompareValue(comparer compare.ValueComparer) *compareValue {
	return &compareValue{
		comparer: comparer,
	}
}
