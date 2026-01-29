// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
)

var _ QueryResultCheck = expectLengthAtLeast{}

type expectLengthAtLeast struct {
	resourceAddress string
	check           int
}

// CheckQuery implements the query check logic.
func (e expectLengthAtLeast) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil {
		resp.Error = fmt.Errorf("no completed query information available")
		return
	}

	if req.QuerySummary.Total < e.check {
		resp.Error = fmt.Errorf("Query result of at least length %v - expected but got %v.", e.check, req.QuerySummary.Total)
		return
	}
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
