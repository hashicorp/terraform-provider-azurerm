// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &mapsAccountPoller{}

type mapsAccountPoller struct {
	client *accounts.AccountsClient
	id     accounts.AccountId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewMapsAccountPoller(client *accounts.AccountsClient, id accounts.AccountId) *mapsAccountPoller {
	return &mapsAccountPoller{
		client: client,
		id:     id,
	}
}

func (p mapsAccountPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ProvisioningState == nil {
		return nil, fmt.Errorf("polling for %s: `provisioningState` was empty", p.id)
	}

	if *resp.Model.Properties.ProvisioningState == string(pollingSuccess.Status) {
		return &pollingSuccess, nil
	}

	return &pollingInProgress, nil
}
