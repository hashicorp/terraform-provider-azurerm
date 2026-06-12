// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &automationWatcherUpdatePoller{}

type automationWatcherUpdatePoller struct {
	client *watcher.WatcherClient
	id     watcher.WatcherId
}

var (
	pollingSuccess = pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		PollInterval: 1 * time.Minute,
		Status:       pollers.PollingStatusInProgress,
	}
	pollingStatusUnknown = pollers.PollResult{
		Status: pollers.PollingStatusUnknown,
	}
)

func NewAutomationWatcherUpdatePoller(client *watcher.WatcherClient, id watcher.WatcherId) *automationWatcherUpdatePoller {
	return &automationWatcherUpdatePoller{
		client: client,
		id:     id,
	}
}

func (p automationWatcherUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			if status := properties.Status; status != nil {
				switch *status {
				case "New":
					return &pollingInProgress, nil
				case "Running", "Stopped", "Suspended":
					pollingSuccess.HttpResponse = &client.Response{
						OData:    resp.OData,
						Response: resp.HttpResponse,
					}
					return &pollingSuccess, nil
				}
			}
		}
	}

	return &pollingStatusUnknown, fmt.Errorf("polling %s update: no expected status is returned", p.id)
}
