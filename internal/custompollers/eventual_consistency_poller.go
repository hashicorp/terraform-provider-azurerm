// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &eventualConsistencyPoller{}

// EventualConsistencyPollerOptions allow for customizing the eventualConsistencyPoller polling logic
type EventualConsistencyPollerOptions struct {
	// Interval specifies the interval until the poller should be called again.
	Interval time.Duration
	// TargetStatusCode specifies the optional target HTTP status code, omitting this treats any non-error response as a success
	TargetStatusCode *int // e.g. http.StatusOK
	// RetryErrorStatusCodes specifies the error HTTP status codes that should be retried
	RetryErrorStatusCodes []int // e.g. http.StatusNotFound, http.StatusForbidden
}

// DefaultCreationEventualConsistencyPollerOptions returns EventualConsistencyPollerOptions that are commonly used
// to account for eventual consistency scenarios where Azure may return a 404 for a period of time after creation
func DefaultCreationEventualConsistencyPollerOptions() *EventualConsistencyPollerOptions {
	return &EventualConsistencyPollerOptions{
		Interval:              10 * time.Second,
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	}
}

// DefaultDeletionEventualConsistencyPollerOptions returns EventualConsistencyPollerOptions that are commonly used
// to account for eventual consistency scenarios where Azure may return continue to return successes (e.g. http.StatusOK) for a period of time after deletion
func DefaultDeletionEventualConsistencyPollerOptions() *EventualConsistencyPollerOptions {
	return &EventualConsistencyPollerOptions{
		Interval:         10 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	}
}

// eventualConsistencyPoller is a poller meant for resources that need to account for eventual consistency
// by checking for n number of successful responses, or n number of responses matching a specific target status code.
type eventualConsistencyPoller struct {
	fn      func(pollerCtx context.Context) (*http.Response, error)
	options EventualConsistencyPollerOptions

	remainingCount int
	targetCount    int
}

// NewEventualConsistencyPoller returns a new pollers.Poller with an eventualConsistencyPoller pollers.pollerType.
// If the provided EventualConsistencyPollerOptions struct is nil, the default creation poller options are used.
func NewEventualConsistencyPoller(targetCount int, fn func(pollerCtx context.Context) (*http.Response, error), o *EventualConsistencyPollerOptions) pollers.Poller {
	if o == nil {
		o = DefaultCreationEventualConsistencyPollerOptions()
	}

	pollerType := &eventualConsistencyPoller{
		fn:      fn,
		options: *o,

		remainingCount: targetCount,
		targetCount:    targetCount,
	}

	return pollers.NewPoller(pollerType, pollerType.options.Interval, pollers.DefaultNumberOfDroppedConnectionsToAllow)
}

func (p *eventualConsistencyPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.fn(ctx)
	if err != nil {
		// In certain cases, we actually want to poll for a specific error status code, such as checking for consecutive 404s on deletion
		if p.options.TargetStatusCode != nil {
			if response.WasStatusCode(resp, *p.options.TargetStatusCode) {
				return p.decrement(), nil
			}
		}

		if response.WasStatusCodes(resp, p.options.RetryErrorStatusCodes...) {
			return p.reset(), nil
		}

		return &pollers.PollResult{
			Status:       pollers.PollingStatusFailed,
			PollInterval: p.options.Interval,
		}, pollers.PollingFailedError{Message: err.Error()}
	}

	// If a target status code is set, we'll check this matches and consider any non-matching code as a failure and reset.
	// otherwise, any non-error code will be considered a success.
	if p.options.TargetStatusCode != nil {
		if response.WasStatusCode(resp, *p.options.TargetStatusCode) {
			return p.decrement(), nil
		}
		return p.reset(), nil
	}

	return p.decrement(), nil
}

func (p *eventualConsistencyPoller) check() *pollers.PollResult {
	if p.remainingCount > 0 {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: p.options.Interval,
		}
	}

	return &pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}
}

func (p *eventualConsistencyPoller) decrement() *pollers.PollResult {
	p.remainingCount--
	return p.check()
}

func (p *eventualConsistencyPoller) reset() *pollers.PollResult {
	p.remainingCount = p.targetCount
	return p.check()
}
