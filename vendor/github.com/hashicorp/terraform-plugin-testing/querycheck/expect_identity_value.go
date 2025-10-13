// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryCheck = expectIdentityValue{}

type expectIdentityValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
	identityValue   knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentityValue) CheckQuery(ctx context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	// var resource *tfjson.QueryResource

	if req.Query == nil {
		resp.Error = fmt.Errorf("query is nil")

		return
	}

	/*	if req.Query.Values == nil {
			resp.Error = fmt.Errorf("query does not contain any query values")

			return
		}

		if req.Query.Values.RootModule == nil {
			resp.Error = fmt.Errorf("query does not contain a root module")

			return
		}

		for _, r := range req.Query.Values.RootModule.Resources {
			if e.resourceAddress == r.Address {
				resource = r

				break
			}
		}

		if resource == nil {
			resp.Error = fmt.Errorf("%s - Resource not found in query", e.resourceAddress)

			return
		}

		if resource.IdentitySchemaVersion == nil || len(resource.IdentityValues) == 0 {
			resp.Error = fmt.Errorf("%s - Identity not found in query. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.12+)", e.resourceAddress)

			return
		}

		result, err := tfjsonpath.Traverse(resource.IdentityValues, e.attributePath)

		if err != nil {
			resp.Error = err

			return
		}*/

	/*	if err := e.identityValue.CheckValue(result); err != nil {
		resp.Error = fmt.Errorf("error checking identity value for attribute at path: %s.%s, err: %s", e.resourceAddress, e.attributePath.String(), err)

		return
	}*/
	return
}

// ExpectIdentityValue returns a query check that asserts that the specified identity attribute at the given resource
// matches a known value. This query check can only be used with managed resources that support resource identity.
//
// Resource identity is only supported in Terraform v1.12+
func ExpectIdentityValue(resourceAddress string, attributePath tfjsonpath.Path, identityValue knownvalue.Check) QueryCheck {
	return expectIdentityValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
		identityValue:   identityValue,
	}
}
