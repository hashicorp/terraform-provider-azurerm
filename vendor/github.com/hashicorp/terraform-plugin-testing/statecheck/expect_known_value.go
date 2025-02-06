// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Resource State Check
var _ StateCheck = expectKnownValue{}

type expectKnownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	knownValue      knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectKnownValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
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

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.knownValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)

		return
	}
}

// ExpectKnownValue returns a state check that asserts that the specified attribute at the given resource
// has a known type and value.
func ExpectKnownValue(resourceAddress string, attributePath tfjsonpath.Path, knownValue knownvalue.Check) StateCheck {
	return expectKnownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		knownValue:      knownValue,
	}
}
