// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pollers

import "context"

type pollKey int

const (
	skipPollingDelayKey pollKey = iota
)

// WithSkipPollingDelay returns a new context with the skip polling delay flag set.
// This is used to signal to PollUntilDone that it should not wait between polling attempts.
func WithSkipPollingDelay(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipPollingDelayKey, true)
}

// ShouldSkipPollingDelay returns true if the context has the skip polling delay flag set.
func ShouldSkipPollingDelay(ctx context.Context) bool {
	if v, ok := ctx.Value(skipPollingDelayKey).(bool); ok {
		return v
	}
	return false
}
