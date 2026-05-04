// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statecheck

import (
	"context"

	tfjson "github.com/hashicorp/terraform-json"
)

// StateCheck defines an interface for implementing test logic that checks a state file and then returns an error
// if the state file does not match what is expected.
type StateCheck interface {
	// CheckState should perform the state check.
	CheckState(context.Context, CheckStateRequest, *CheckStateResponse)
}

// CheckStateRequest is a request for an invoke of the CheckState function.
type CheckStateRequest struct {
	// State represents a parsed state file, retrieved via the `terraform show -json` command.
	State *tfjson.State
}

// CheckStateResponse is a response to an invoke of the CheckState function.
type CheckStateResponse struct {
	// Error is used to report the failure of a state check assertion and is combined with other StateCheck errors
	// to be reported as a test failure.
	Error error
}
