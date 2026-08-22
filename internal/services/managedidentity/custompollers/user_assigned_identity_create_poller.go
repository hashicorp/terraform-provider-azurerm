// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &userAssignedIdentityCreatePoller{}

const defaultSuccessCount = 5

var (
	pollingSuccess = &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
	pollingFailed = &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusFailed,
	}
)

type userAssignedIdentityCreatePoller struct {
	client       *identities.IdentitiesClient
	id           commonids.UserAssignedIdentityId
	successCount int
}

func NewUserAssignedIdentityCreatePoller(client *identities.IdentitiesClient, id commonids.UserAssignedIdentityId) *userAssignedIdentityCreatePoller {
	return &userAssignedIdentityCreatePoller{
		client:       client,
		id:           id,
		successCount: defaultSuccessCount,
	}
}

func (p *userAssignedIdentityCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.UserAssignedIdentitiesGet(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			p.successCount = defaultSuccessCount
			return pollingInProgress, nil
		}
		return pollingFailed, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if p.successCount > 1 {
		p.successCount--
		return pollingInProgress, nil
	}

	return pollingSuccess, nil
}
