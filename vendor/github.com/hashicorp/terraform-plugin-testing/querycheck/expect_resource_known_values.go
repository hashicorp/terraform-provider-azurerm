// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ QueryResultCheck = expectResourceKnownValues{}
var _ QueryResultCheckWithFilters = expectResourceKnownValues{}

type expectResourceKnownValues struct {
	listResourceAddress string
	filter              queryfilter.QueryFilter
	knownValueChecks    []KnownValueCheck
}

func (e expectResourceKnownValues) QueryFilters(ctx context.Context) []queryfilter.QueryFilter {
	if e.filter == nil {
		return []queryfilter.QueryFilter{}
	}

	return []queryfilter.QueryFilter{
		e.filter,
	}
}

func (e expectResourceKnownValues) CheckQuery(_ context.Context, req CheckQueryRequest, resp *CheckQueryResponse) {
	listRes := make([]tfjson.ListResourceFoundData, 0)
	var diags []error
	for _, res := range req.Query {
		if e.listResourceAddress == strings.TrimPrefix(res.Address, "list.") {
			listRes = append(listRes, res)
		}
	}

	if len(listRes) == 0 {
		resp.Error = fmt.Errorf("%s - no query results found after filtering", e.listResourceAddress)
		return
	}

	if len(listRes) > 1 {
		resp.Error = fmt.Errorf("%s - more than 1 query result found after filtering", e.listResourceAddress)
		return
	}

	res := listRes[0]

	if res.ResourceObject == nil {
		resp.Error = fmt.Errorf("%s - no resource object was returned, ensure `include_resource` has been set to `true` in the list resource config`", e.listResourceAddress)
		return
	}

	for _, c := range e.knownValueChecks {
		resource, err := tfjsonpath.Traverse(res.ResourceObject, c.Path)
		if err != nil {
			resp.Error = err
			return
		}

		if err := c.KnownValue.CheckValue(resource); err != nil {
			diags = append(diags, fmt.Errorf("error checking value for attribute at path: %s for resource with identity %s, err: %s", c.Path.String(), e.filter, err))
		}
	}

	if diags != nil {
		var diagsStr string
		for _, diag := range diags {
			diagsStr += diag.Error() + "; "
		}
		resp.Error = fmt.Errorf("the following errors were found while checking values: %s", diagsStr)
		return
	}
}

// ExpectResourceKnownValues returns a query check which asserts that a resource object identified by a query filter
// passes the given query checks.
//
// This query check can only be used with managed resources that support resource identity and query. Query is only supported in Terraform v1.14+
func ExpectResourceKnownValues(listResourceAddress string, filter queryfilter.QueryFilter, knownValues []KnownValueCheck) QueryResultCheck {
	return expectResourceKnownValues{
		listResourceAddress: listResourceAddress,
		filter:              filter,
		knownValueChecks:    knownValues,
	}
}
