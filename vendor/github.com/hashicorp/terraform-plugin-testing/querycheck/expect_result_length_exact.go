// Copyright IBM Corp. 2014, 2025
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
	if req.QuerySummary == nil && req.QuerySummaries == nil {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	if strings.HasPrefix(e.resourceAddress, "list.") {
		e.resourceAddress = strings.TrimPrefix(e.resourceAddress, "list.")
	}

	total := 0
	for _, summary := range req.QuerySummaries {
		// To support query tests where for_each is used to construct the list blocks dynamically (e.g. with child resources) we allow
		// specifying a trailing '[*]' to indicate that we should be looking for multiple summaries

		if !strings.HasSuffix(e.resourceAddress, "[*]") {
			if strings.EqualFold(strings.TrimPrefix(summary.Address, "list."), e.resourceAddress) {
				if e.check != summary.Total {
					resp.Error = fmt.Errorf("number of found resources for %s - expected %v but got %v", e.resourceAddress, e.check, summary.Total)
				}
				// It's been found and checked, we can exit
				return
			}
		} else {
			// when using for_each summary.Address returns as `list.{resource_name}.{resource_label}[{each.key}]`
			if strings.HasPrefix(strings.TrimPrefix(summary.Address, "list."), strings.TrimSuffix(e.resourceAddress, "*]")) {
				total = total + summary.Total
			}
		}
	}

	if total != e.check {
		resp.Error = fmt.Errorf("number of found resources for %s - expected %v but got %v", e.resourceAddress, e.check, total)
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
