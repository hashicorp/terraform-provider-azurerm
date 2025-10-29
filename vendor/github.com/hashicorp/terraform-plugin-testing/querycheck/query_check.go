// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// QueryCheck defines an interface for implementing test logic that checks a query file and then returns an error
// if the query file does not match what is expected.
type QueryCheck interface {
	// CheckQuery should perform the query check.
	CheckQuery(context.Context, CheckQueryRequest, *CheckQueryResponse)
}

// CheckQueryRequest is a request for an invoke of the CheckQuery function.
type CheckQueryRequest struct {
	// Query represents a parsed query file, retrieved via the `terraform show -json` command.
	Query *[]tfjson.LogMsg
}

// CheckQueryResponse is a response to an invoke of the CheckQuery function.
type CheckQueryResponse struct {
	// Error is used to report the failure of a query check assertion and is combined with other QueryCheck errors
	// to be reported as a test failure.
	Error error
}
