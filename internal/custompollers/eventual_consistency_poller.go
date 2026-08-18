package custompollers

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &eventualConsistencyPoller{}

// eventualConsistencyPoller is a general use poller meant to be used with resources that may return a 404 for a period of time shortly after creation.
// This is the standard replacement for `StateChangeConf` configurations using the `ContinuousTargetOccurence` field that only check for existence.
type eventualConsistencyPoller struct {
	fn       func() (*http.Response, error)
	interval time.Duration

	remainingCount int
	successCount   int
}

func NewEventualConsistencyPoller(interval time.Duration, successCount int, fn func() (*http.Response, error)) pollers.Poller {
	pollerType := &eventualConsistencyPoller{
		fn:             fn,
		interval:       interval,
		remainingCount: successCount,
		successCount:   successCount,
	}

	return pollers.NewPoller(pollerType, pollerType.interval, pollers.DefaultNumberOfDroppedConnectionsToAllow)
}

func (p *eventualConsistencyPoller) Poll(_ context.Context) (*pollers.PollResult, error) {
	resp, err := p.fn()
	if err != nil {
		if response.WasNotFound(resp) {
			p.remainingCount = p.successCount
			return &pollers.PollResult{
				Status:       pollers.PollingStatusInProgress,
				PollInterval: p.interval,
			}, nil
		}

		return &pollers.PollResult{
			Status:       pollers.PollingStatusFailed,
			PollInterval: p.interval,
		}, pollers.PollingFailedError{Message: err.Error()}
	}

	p.remainingCount--
	if p.remainingCount > 0 {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: p.interval,
		}, nil
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: p.interval,
	}, nil
}
