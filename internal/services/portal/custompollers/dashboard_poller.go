// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dashboardPoller{}

type dashboardPoller struct {
	client *dashboard.DashboardClient
	id     dashboard.DashboardId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

// NewDashboardPoller returns a poller which waits for a newly created Dashboard to become visible to `Get`.
//
// `CreateOrUpdate` for this API is not a Long Running Operation and can return before the Dashboard has been
// propagated, in which case an immediate `Get` returns a 404.
func NewDashboardPoller(client *dashboard.DashboardClient, id dashboard.DashboardId) *dashboardPoller {
	return &dashboardPoller{
		client: client,
		id:     id,
	}
}

func (p dashboardPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingInProgress, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollingSuccess, nil
}
