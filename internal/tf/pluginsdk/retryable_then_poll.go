// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"context"
	"net/http"
	"time"
)

// Poller is the minimal interface required for long-running operations.
// It is intentionally compatible with pollers returned by go-azure-sdk.
type Poller interface {
	PollUntilDone(ctx context.Context) error
}

// RetryableThenPoll performs an operation which may return a poller, retrying only when
// isRetryable returns true for the operation error.
//
// When the operation succeeds, any returned poller is executed to completion.
func RetryableThenPoll(ctx context.Context, timeout time.Duration, operation func(context.Context) (Poller, *http.Response, error), isRetryable func(*http.Response, error) bool) error {
	return Retry(timeout, func() *RetryError {
		poller, httpResponse, err := operation(ctx)
		if err != nil {
			if isRetryable != nil && isRetryable(httpResponse, err) {
				return RetryableError(err)
			}
			return NonRetryableError(err)
		}

		if poller == nil {
			return nil
		}

		if err := poller.PollUntilDone(ctx); err != nil {
			return NonRetryableError(err)
		}

		return nil
	})
}
