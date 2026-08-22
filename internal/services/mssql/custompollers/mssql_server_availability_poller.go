// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &MsSqlServerAvailabilityPoller{}

type MsSqlServerAvailabilityPoller struct {
	client       *servers.ServersClient
	id           commonids.SqlServerId
	successCount int
}

const msSqlServerAvailabilitySuccessCount = 3

func NewMsSqlServerAvailabilityPoller(client *servers.ServersClient, id commonids.SqlServerId) *MsSqlServerAvailabilityPoller {
	return &MsSqlServerAvailabilityPoller{
		client:       client,
		id:           id,
		successCount: msSqlServerAvailabilitySuccessCount,
	}
}

func (p *MsSqlServerAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, servers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			p.successCount = msSqlServerAvailabilitySuccessCount
			return &pollers.PollResult{
				Status:       pollers.PollingStatusInProgress,
				PollInterval: 5 * time.Second,
			}, nil
		}

		return &pollers.PollResult{
			Status:       pollers.PollingStatusFailed,
			PollInterval: 5 * time.Second,
		}, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if p.successCount > 1 {
		p.successCount--
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: 5 * time.Second,
		}, nil
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 5 * time.Second,
	}, nil
}
