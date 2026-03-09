// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &ResourceProviderExistencePoller{}

// NewResourceProviderExistencePoller - Checks for RPs in a "Registered" state, retrying n times
func NewResourceProviderExistencePoller(client *providers.ProvidersClient, id providers.SubscriptionProviderId, maxRetries int) *ResourceProviderExistencePoller {
	return &ResourceProviderExistencePoller{
		client:     client,
		id:         id,
		maxRetries: maxRetries,
	}
}

type ResourceProviderExistencePoller struct {
	client     *providers.ProvidersClient
	id         providers.SubscriptionProviderId
	maxRetries int
}

func (p *ResourceProviderExistencePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.maxRetries == 0 {
		return &pollers.PollResult{
			Status: pollers.PollingStatusCancelled,
		}, pollers.PollingCancelledError{Message: "maximum retries exceeded"}
	}

	resp, err := p.client.Get(ctx, p.id, providers.DefaultGetOperationOptions())
	if err != nil {
		return nil, pollers.PollingFailedError{Message: fmt.Sprintf("retrieving %s: %+v", p.id, err)}
	}

	if resp.Model == nil {
		return nil, pollers.PollingFailedError{Message: fmt.Sprintf("retrieving %s: `model` was nil", p.id)}
	}

	if strings.EqualFold(pointer.From(resp.Model.RegistrationState), "Registered") {
		return &pollers.PollResult{
			Status: pollers.PollingStatusSucceeded,
		}, nil
	} else {
		p.maxRetries -= 1
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
