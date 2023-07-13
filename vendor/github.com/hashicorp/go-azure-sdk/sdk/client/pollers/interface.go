// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

// PollerType allows custom pollers to be created to determine when a particular Operation has
// been Completed, Cancelled or Failed.
type PollerType interface {
	// Poll performs a poll to determine whether the Operation has been Completed/Cancelled or Failed.
	// Behaviourally this method should:
	//   1. Perform a single poll, and assume that it will be called again after the PollInterval
	//      defined within the PollResult.
	//   2. When the operation is still in progress, a PollResult should be returned with the Status
	//      `InProgress` and the next PollInterval.
	//   3. When the operation is Completed, a PollResult should be returned with the Status
	//      `Succeeded`.
	//   4. When the operation is Cancelled a PollingCancelledError should be returned.
	//   5. When the operation Fails a PollingFailedError should be returned.
	Poll(ctx context.Context) (*PollResult, error)
}

type PollResult struct {
	// HttpResponse is a copy of the HttpResponse returned from the API in the last request.
	HttpResponse *client.Response

	// PollInterval specifies the interval until this poller should be called again.
	PollInterval time.Duration

	// Status specifies the polling status of this resource at the time of the last request.
	Status PollingStatus
}

type PollingStatus string

const (
	// PollingStatusCancelled states that the resource change has been cancelled.
	PollingStatusCancelled PollingStatus = "Cancelled"

	// PollingStatusFailed states that the resource change has Failed.
	PollingStatusFailed PollingStatus = "Failed"

	// PollingStatusInProgress states that the resource change is still occurring/in-progress.
	PollingStatusInProgress PollingStatus = "InProgress"

	// PollingStatusSucceeded states that the resource change was successful.
	PollingStatusSucceeded PollingStatus = "Succeeded"

	// PollingStatusUnknown states that the resource change state is unknown/unexpected.
	PollingStatusUnknown PollingStatus = "Unknown"
)

var _ error = PollingCancelledError{}

// PollingCancelledError defines the that the resource change was cancelled (for example, due to a timeout).
type PollingCancelledError struct {
	// HttpResponse is a copy of the HttpResponse returned from the API in the last request.
	HttpResponse *client.Response

	// Message is a custom error message containing more details
	Message string
}

func (e PollingCancelledError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("polling was cancelled: %+v", e.Message)
	}

	return fmt.Sprintf("polling was cancelled")
}

var _ error = PollingDroppedConnectionError{}

// PollingDroppedConnectionError defines there was a dropped connection when polling for the status.
type PollingDroppedConnectionError struct {
	// Message is a custom error message containing more details
	Message string
}

func (e PollingDroppedConnectionError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("experienced a dropped connection when polling: %+v", e.Message)
	}

	return fmt.Sprintf("experienced a dropped connection when polling")
}

var _ error = PollingFailedError{}

// PollingFailedError states that the resource change failed (for example due to a lack of capacity).
type PollingFailedError struct {
	// HttpResponse is a copy of the HttpResponse returned from the API in the last request.
	HttpResponse *client.Response

	// Message is a custom error message containing more details
	Message string
}

func (e PollingFailedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("polling failed: %+v", e.Message)
	}

	return fmt.Sprintf("polling failed")
}
