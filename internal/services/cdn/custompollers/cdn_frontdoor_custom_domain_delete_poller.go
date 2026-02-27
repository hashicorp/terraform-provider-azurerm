// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorCustomDomainDeletePoller{}

type frontDoorCustomDomainDeletePoller struct {
	client *afdcustomdomains.AFDCustomDomainsClient
	id     afdcustomdomains.CustomDomainId
}

var (
	frontDoorCustomDomainDeleteSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	frontDoorCustomDomainDeleteInProgress = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewFrontDoorCustomDomainDeletePoller(client *afdcustomdomains.AFDCustomDomainsClient, id afdcustomdomains.CustomDomainId) pollers.PollerType {
	return &frontDoorCustomDomainDeletePoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorCustomDomainDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &frontDoorCustomDomainDeleteSuccess, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &frontDoorCustomDomainDeleteInProgress, nil
}
