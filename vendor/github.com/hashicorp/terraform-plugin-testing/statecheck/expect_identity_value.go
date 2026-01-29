// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ StateCheck = expectIdentityValue{}

type expectIdentityValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	identityValue   knownvalue.Check
}

// CheckState implements the state check logic.
func (e expectIdentityValue) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
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

	if resource.IdentitySchemaVersion == nil || len(resource.IdentityValues) == 0 {
		resp.Error = fmt.Errorf("%s - Identity not found in state. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.12+)", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.IdentityValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	if err := e.identityValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking identity value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)

		return
	}
}

// ExpectIdentityValue returns a state check that asserts that the specified identity attribute at the given resource
// matches a known value. This state check can only be used with managed resources that support resource identity.
//
// Resource identity is only supported in Terraform v1.12+
func ExpectIdentityValue(resourceAddress string, attributePath tfjsonpath.Path, identityValue knownvalue.Check) StateCheck {
	return expectIdentityValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		identityValue:   identityValue,
	}
}
