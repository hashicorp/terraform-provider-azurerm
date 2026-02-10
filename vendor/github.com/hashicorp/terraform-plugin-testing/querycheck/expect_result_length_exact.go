// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
)

var _ QueryResultCheck = expectLength{}

type expectLength struct {
	resourceAddress string
	check           int
}

// CheckQuery implements the query check logic.
func (e expectLength) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	if e.check != req.QuerySummary.Total {
		resp.Error = fmt.Errorf("number of found resources %v - expected but got %v.", e.check, req.QuerySummary.Total)
		return
	}
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
