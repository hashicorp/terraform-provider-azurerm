// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type dataFactoryEventSubscriptionStatusPoller struct {
	client datafactory.TriggersClient
	id     parse.TriggerId
}

var _ pollers.PollerType = &dataFactoryEventSubscriptionStatusPoller{}

func NewDataFactoryEventSubscriptionPoller(client datafactory.TriggersClient, id parse.TriggerId) pollers.PollerType {
	return &dataFactoryEventSubscriptionStatusPoller{
		client: client,
		id:     id,
	}
}

func (p dataFactoryEventSubscriptionStatusPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result, err := p.client.GetEventSubscriptionStatus(ctx, p.id.ResourceGroup, p.id.FactoryName, p.id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving event subscription status for %s: %+v", p.id, err)
	}

	if result.Status == datafactory.EventSubscriptionStatusEnabled {
		return &pollers.PollResult{
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	if result.Status != datafactory.EventSubscriptionStatusProvisioning {
		return nil, fmt.Errorf("event subscription for %s reached unexpected status `%s`", p.id, result.Status)
	}

	return &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
