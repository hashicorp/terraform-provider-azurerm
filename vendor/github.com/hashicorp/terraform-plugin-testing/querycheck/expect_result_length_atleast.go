// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"strings"
)

var _ QueryResultCheck = expectLengthAtLeast{}

type expectLengthAtLeast struct {
	resourceAddress string
	check           int
}

// CheckQuery implements the query check logic.
func (e expectLengthAtLeast) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil && len(req.QuerySummaries) == 0 {
		resp.Error = fmt.Errorf("no completed query information available")
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
			if summary.Total < e.check {
				resp.Error = fmt.Errorf("Query result of at least length %v - expected but got %v.", e.check, summary.Total)
				return
			}
			return
		}
	}

	resp.Error = fmt.Errorf("the list block %s was not found in the query results", e.resourceAddress)
}

// ExpectLengthAtLeast returns a query check that asserts that the length of the query result is at least the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectLengthAtLeast(resourceAddress string, length int) QueryResultCheck {
	return expectLengthAtLeast{
		resourceAddress: resourceAddress,
		check:           length,
	}
}
