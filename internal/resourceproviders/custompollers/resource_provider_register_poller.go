// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

const CONTINUES_TARGET_OCCURENCE = 5

var _ pollers.PollerType = &resourceProviderRegistrationPoller{}

func NewResourceProviderRegistrationPoller(client *providers.ProvidersClient, id providers.SubscriptionProviderId) *resourceProviderRegistrationPoller {
	return &resourceProviderRegistrationPoller{
		client: client,
		id:     id,
	}
}

type resourceProviderRegistrationPoller struct {
	client                   *providers.ProvidersClient
	id                       providers.SubscriptionProviderId
	continuesTargetOccurence int
}

func (p *resourceProviderRegistrationPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, providers.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	registrationState := ""
	if model := resp.Model; model != nil && model.RegistrationState != nil {
		registrationState = *model.RegistrationState
	}

	if strings.EqualFold(registrationState, "Registered") {
		if p.continuesTargetOccurence == CONTINUES_TARGET_OCCURENCE {
			return &pollers.PollResult{
				Status:       pollers.PollingStatusSucceeded,
				PollInterval: 10 * time.Second,
			}, nil
		}
		p.continuesTargetOccurence += 1
	} else {
		p.continuesTargetOccurence = 0
	}

	// Processing
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
