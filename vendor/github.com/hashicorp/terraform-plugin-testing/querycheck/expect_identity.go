// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

var _ QueryResultCheck = expectIdentity{}

type expectIdentity struct {
	listResourceAddress string
	check               map[string]knownvalue.Check
}

// CheckQuery implements the query check logic.
func (e expectIdentity) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	for _, res := range req.Query {
		var errCollection []error

		if e.listResourceAddress != strings.TrimPrefix(res.Address, "list.") {
			continue
		}

		if len(res.Identity) != len(e.check) {
			deltaMsg := ""
			if len(res.Identity) > len(e.check) {
				deltaMsg = statecheck.CreateDeltaString(res.Identity, e.check, "actual identity has extra attribute(s): ")
			} else {
				deltaMsg = statecheck.CreateDeltaString(e.check, res.Identity, "actual identity is missing attribute(s): ")
			}

			resp.Error = fmt.Errorf("%s - Expected %d attribute(s) in the actual identity object, got %d attribute(s): %s", e.listResourceAddress, len(e.check), len(res.Identity), deltaMsg)
			return
		}

		var keys []string

		for k := range e.check {
			keys = append(keys, k)
		}

		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, k := range keys {
			actualIdentityVal, ok := res.Identity[k]

			if !ok {
				resp.Error = fmt.Errorf("%s - missing attribute %q in actual identity object", e.listResourceAddress, k)
				return
			}

			if err := e.check[k].CheckValue(actualIdentityVal); err != nil {
				errCollection = append(errCollection, fmt.Errorf("%s - %q identity attribute: %s", e.listResourceAddress, k, err))
			}
		}

		if errCollection == nil {
			return
		}
	}

	var errCollection []error
	errCollection = append(errCollection, fmt.Errorf("an identity with the following attributes was not found"))

	var keys []string

	for k := range e.check {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// wrap errors for each check
	for _, attr := range keys {
		check := e.check[attr]
		errCollection = append(errCollection, fmt.Errorf("attribute %q: %s", attr, check))
	}

	errCollection = append(errCollection, fmt.Errorf("address: %s\n", e.listResourceAddress))
	resp.Error = errors.Join(errCollection...)
}

// ExpectIdentity returns a query check that asserts that the identity at the given resource matches a known object, where each
// map key represents an identity attribute name. The identity in query must exactly match the given object.
//
// This query check can only be used with managed resources that support resource identity and query. Query is only supported in Terraform v1.14+
func ExpectIdentity(resourceAddress string, identity map[string]knownvalue.Check) QueryResultCheck {
	return expectIdentity{
		listResourceAddress: resourceAddress,
		check:               identity,
	}
}
