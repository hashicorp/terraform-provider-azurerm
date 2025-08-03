// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ statecheck.StateCheck = expectStateContainsIdentityValueAtPath{}

type expectStateContainsIdentityValueAtPath struct {
	resourceAddress  string
	identityAttrPath tfjsonpath.Path
	stateAttrPath    tfjsonpath.Path
}

func (e expectStateContainsIdentityValueAtPath) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
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

	identityString, iOk := identityResult.(string)
	stateString, sOk := stateResult.(string)
	if !iOk || !sOk {
		resp.Error = errors.New("ExpectStateContainsIdentityValueAtPath only works with string state + identity attributes")

		return
	}

	if !strings.Contains(stateString, identityString) {
		resp.Error = fmt.Errorf(
			"expected state (%[1]s.%[2]s) to contain identity value (%[1]s.%[3]s): identity value: %[4]v, state value: %[5]v",
			e.resourceAddress,
			e.stateAttrPath.String(),
			e.identityAttrPath.String(),
			identityResult,
			stateResult,
		)

		return
	}
}

func ExpectStateContainsIdentityValueAtPath(resourceAddress string, identityAttrPath, stateAttrPath tfjsonpath.Path) statecheck.StateCheck {
	return expectStateContainsIdentityValueAtPath{
		resourceAddress:  resourceAddress,
		identityAttrPath: identityAttrPath,
		stateAttrPath:    stateAttrPath,
	}
}
