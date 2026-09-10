// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorCustomDomainDeletePoller{}

type frontDoorCustomDomainDeletePoller struct {
	client *afddomains.AFDDomainsClient
	id     afddomains.CustomDomainId
}

func NewFrontDoorCustomDomainDeletePoller(client *afddomains.AFDDomainsClient, id afddomains.CustomDomainId) pollers.PollerType {
	return &frontDoorCustomDomainDeletePoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorCustomDomainDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.AFDCustomDomainsGet(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusSucceeded,
			}, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
