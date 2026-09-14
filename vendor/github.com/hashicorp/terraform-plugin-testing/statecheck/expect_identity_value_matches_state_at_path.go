// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ StateCheck = expectIdentityValueMatchesStateAtPath{}

type expectIdentityValueMatchesStateAtPath struct {
	resourceAddress  string
	identityAttrPath tfjsonpath.Path
	stateAttrPath    tfjsonpath.Path
}

// CheckState implements the state check logic.
func (e expectIdentityValueMatchesStateAtPath) CheckState(ctx context.Context, req CheckStateRequest, resp *CheckStateResponse) {
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

	identityResult, err := tfjsonpath.Traverse(resource.IdentityValues, e.identityAttrPath)

	if err != nil {
		resp.Error = err

		return
	}

	stateResult, err := tfjsonpath.Traverse(resource.AttributeValues, e.stateAttrPath)

	if err != nil {
		resp.Error = err

		return
	}

	if !reflect.DeepEqual(identityResult, stateResult) {
		resp.Error = fmt.Errorf(
			"expected identity (%[1]s.%[2]s) and state value (%[1]s.%[3]s) to match, but they differ: identity value: %[4]v, state value: %[5]v",
			e.resourceAddress,
			e.identityAttrPath.String(),
			e.stateAttrPath.String(),
			identityResult,
			stateResult,
		)

		return
	}
}

// ExpectIdentityValueMatchesStateAtPath returns a state check that asserts that the specified identity attribute at the given resource
// matches the specified attribute in state. This is useful when an identity attribute is in sync with a state attribute of a different path.
//
// This state check can only be used with managed resources that support resource identity. Resource identity is only supported in Terraform v1.12+
func ExpectIdentityValueMatchesStateAtPath(resourceAddress string, identityAttrPath, stateAttrPath tfjsonpath.Path) StateCheck {
	return expectIdentityValueMatchesStateAtPath{
		resourceAddress:  resourceAddress,
		identityAttrPath: identityAttrPath,
		stateAttrPath:    stateAttrPath,
	}
}
