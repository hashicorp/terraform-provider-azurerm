// Copyright (c) HashiCorp, Inc.
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
	//if req.QuerySummary == nil {
	//	resp.Error = fmt.Errorf("no query summary information available")
	//	return
	//}
	//
	//if e.check != req.QuerySummary.Total {
	//	resp.Error = fmt.Errorf("number of found resources %v - expected but got %v.", e.check, req.QuerySummary.Total)
	//	return
	//}

	found := 0
	for _, resource := range req.Query {
		// when using for_each, resource.Address returns as `list.{resource_name}.{resource_label}[{each.key}]`
		// below for testing/proof of concept, would need more appropriate parsing to ensure resources aren't counted unintentionally
		// e.g.
		// - `list.azurerm_example.test` + `list.azurerm_example.test2` would both match prefix `list.azurerm_example.test`
		if strings.HasPrefix(resource.Address, e.resourceAddress) {
			found++
		}
	}

	if found != e.check {
		resp.Error = fmt.Errorf("found resources for %s - %d expected but got %d", e.resourceAddress, e.check, found)
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
