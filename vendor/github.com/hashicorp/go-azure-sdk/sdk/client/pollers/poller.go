// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pollers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

const DefaultNumberOfDroppedConnectionsToAllow = 3

type Poller struct {
	// initialDelayDuration specifies the duration of the initial delay when polling
	// this is also used for retries should a `latestResponse` not be available, for
	// example when a connection is dropped.
	initialDelayDuration time.Duration

	// latestError contains the error returned from the latest poll.
	latestError error

	// latestResponse contains the polling status from the latest response.
	latestResponse *PollResult

	// maxNumberOfDroppedConnections specifies the maximum number of sequential dropped connections before an error is raised.
	maxNumberOfDroppedConnections int

	// poller is a reference to the PollerType, for example a LongRunningOperationPoller
	// which should be polled to determine the latest state.
	poller PollerType
}

func NewPoller(pollerType PollerType, initialDelayDuration time.Duration, maxNumberOfDroppedConnections int) Poller {
	return Poller{
		initialDelayDuration:          initialDelayDuration,
		maxNumberOfDroppedConnections: maxNumberOfDroppedConnections,
		poller:                        pollerType,
	}
}

// LatestResponse returns the latest HTTP Response returned when polling
func (p *Poller) LatestResponse() *client.Response {
	if p.latestError != nil {
		if v, ok := p.latestError.(PollingCancelledError); ok {
			return v.HttpResponse
		}
		if _, ok := p.latestError.(PollingDroppedConnectionError); ok {
			return nil
		}
		if v, ok := p.latestError.(PollingFailedError); ok {
			return v.HttpResponse
		}

		if p.latestError == context.DeadlineExceeded {
			return nil
		}
	}

	if p.latestResponse == nil {
		return nil
	}

	return p.latestResponse.HttpResponse
}

// LatestStatus returns the latest status returned when polling
func (p *Poller) LatestStatus() PollingStatus {
	if p.latestError != nil {
		if _, ok := p.latestError.(PollingCancelledError); ok {
			return PollingStatusCancelled
		}
		if _, ok := p.latestError.(PollingDroppedConnectionError); ok {
			// we could look to expose a status for this, but we likely wouldn't handle this any differently
			// to it being unknown, so I (@tombuildsstuff) think this is reasonable for now?
			return PollingStatusUnknown
		}
		if _, ok := p.latestError.(PollingFailedError); ok {
			return PollingStatusFailed
		}
		if p.latestError == context.DeadlineExceeded {
			return PollingStatusUnknown
		}
	}

	if p.latestResponse == nil {
		return PollingStatusUnknown
	}

	return p.latestResponse.Status
}

// PollUntilDone polls until the poller determines that the operation has been completed
func (p *Poller) PollUntilDone(ctx context.Context) error {
	if p.poller == nil {
		return fmt.Errorf("internal-error: `poller` was nil`")
	}
	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("internal-error: `ctx` should have a deadline")
	}

	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		connectionDropCounter := 0
		retryDuration := p.initialDelayDuration
		for true {
			// determine the next retry duration / how long to poll for
			if p.latestResponse != nil {
				retryDuration = p.latestResponse.PollInterval
			}
			endTime := time.Now().Add(retryDuration)
			select {
			case <-time.After(time.Until(endTime)):
				{
					break
				}
			}

			p.latestResponse, p.latestError = p.poller.Poll(ctx)

			// first check the connection drop status
			connectionHasBeenDropped := false
			if p.latestResponse == nil && p.latestError == nil {
				// connection drops can either have no response/error (where we have no context)
				connectionHasBeenDropped = true
			} else if _, ok := p.latestError.(PollingDroppedConnectionError); ok {
				// or have an error with more details (e.g. server not found, connection reset etc)
				connectionHasBeenDropped = true
			}
			if connectionHasBeenDropped {
				connectionDropCounter++
				if connectionDropCounter < p.maxNumberOfDroppedConnections {
					continue
				}
				if p.latestResponse == nil && p.latestError == nil {
					// the connection was dropped, but we have no context
					p.latestError = PollingDroppedConnectionError{}
					break
				}
			} else {
				connectionDropCounter = 0
			}

			if p.latestError != nil {
				break
			}

			if response := p.latestResponse; response != nil {
				retryDuration = response.PollInterval

				done := false
				switch response.Status {
				// Cancelled, Dropped Connections and Failed should be raised as errors containing additional info if available

				case PollingStatusCancelled:
					p.latestError = fmt.Errorf("internal-error: a polling status of `Cancelled` should be surfaced as a PollingCancelledError")
					done = true
					break

				case PollingStatusFailed:
					p.latestError = fmt.Errorf("internal-error: a polling status of `Failed` should be surfaced as a PollingFailedError")
					done = true
					break

				case PollingStatusInProgress:
					continue

				case PollingStatusSucceeded:
					done = true
					break

				default:
					p.latestError = fmt.Errorf("internal-error: unimplemented polling status %q", string(response.Status))
					done = true
					break
				}

				if done {
					break
				}
			}
		}
		wait.Done()
	}()

	waitDone := make(chan struct{}, 1)
	go func() {
		wait.Wait()
		waitDone <- struct{}{}
	}()

	select {
	case <-waitDone:
		break
	case <-ctx.Done():
		{
			p.latestResponse = nil
			p.latestError = ctx.Err()
			return p.latestError
		}
	}

	if p.latestError != nil {
		p.latestResponse = nil
	}

	return p.latestError
}
