// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"
)

var _ QueryResultCheck = expectLength{}

type expectLength struct {
	resourceAddress string
	check           int
}

// CheckQuery implements the query check logic.
func (e expectLength) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil && len(req.QuerySummaries) == 0 {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	for _, summary := range req.QuerySummaries {
		address := summary.Address

		// this brings the behaviour of this check in-line with the other query checks where the resource
		// address needs to be provided without the `list.` prefix, but maintains the previous behaviour
		// to not break existing tests that may be using the `list.` prefix in the resource address
		if !strings.HasPrefix(e.resourceAddress, "list.") {
			address = strings.TrimPrefix(summary.Address, "list.")
		}

		if strings.EqualFold(address, e.resourceAddress) {
			if summary.Total != e.check {
				resp.Error = fmt.Errorf("number of found resources %v - expected but got %v.", e.check, summary.Total)
				return
			}
			return
		}
	}

	resp.Error = fmt.Errorf("the list block %s was not found in the query results", e.resourceAddress)
}

// ExpectLength returns a query check that asserts that the length of the query result is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLength(resourceAddress string, length int) QueryResultCheck {
	return expectLength{
		resourceAddress: resourceAddress,
		check:           length,
	}
}
