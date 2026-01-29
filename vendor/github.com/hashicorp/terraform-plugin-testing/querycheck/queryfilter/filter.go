// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package queryfilter

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// QueryFilter defines an interface for implementing declarative filtering logic to apply to query results before
// the results are passed to a query check request.
type QueryFilter interface {
	Filter(context.Context, FilterQueryRequest, *FilterQueryResponse)
}

// FilterQueryRequest is a request to a filter function.
type FilterQueryRequest struct {
	// QueryItem represents a single parsed log message relating to a found resource returned by the `terraform query -json` command.
	QueryItem tfjson.ListResourceFoundData
}

// FilterQueryResponse is a response to a filter function.
type FilterQueryResponse struct {
	// Include indicates whether the QueryItem should be included in CheckQueryRequest.Query
	Include bool

	// Error is used to report the failure of filtering and is combined with other QueryFilter errors
	// to be reported as a test failure.
	Error error
}
