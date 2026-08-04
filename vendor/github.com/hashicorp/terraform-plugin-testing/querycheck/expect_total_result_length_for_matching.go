// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	"regexp"
)

var _ QueryResultCheck = expectTotalLengthForMatching{}

type expectTotalLengthForMatching struct {
	regex *regexp.Regexp
	check int
}

// CheckQuery implements the query check logic.
func (e expectTotalLengthForMatching) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	if req.QuerySummary == nil && len(req.QuerySummaries) == 0 {
		resp.Error = fmt.Errorf("no query summary information available")
		return
	}

	total := 0
	matchFound := false
	for _, summary := range req.QuerySummaries {
		if e.regex.MatchString(summary.Address) {
			total += summary.Total
			matchFound = true
		}
	}

	if !matchFound {
		resp.Error = fmt.Errorf("no list resources matching the provided regex pattern %s were found in the query results", e.regex.String())
		return
	}

	if total != e.check {
		resp.Error = fmt.Errorf("expected total of found resources to be %d, got %d", e.check, total)
	}
}

// ExpectTotalLengthForMatching returns a query check that asserts that the sum of query result lengths
// produced by multiple list blocks is exactly the given value.
//
// This query check can only be used with managed resources that support query. Query is only supported in Terraform v1.14+
func ExpectTotalLengthForMatching(regex *regexp.Regexp, length int) QueryResultCheck {
	return expectTotalLengthForMatching{
		regex: regex,
		check: length,
	}
}
