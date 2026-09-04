// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &subscriptionPricingUpdatePoller{}

type subscriptionPricingUpdatePoller struct {
	client  *pricings.PricingsClient
	id      pricings.PricingId
	payload pricings.Pricing

	// Response contains the result of the `Update` call that succeeded.
	Response pricings.UpdateOperationResponse
}

func NewSubscriptionPricingUpdatePoller(client *pricings.PricingsClient, id pricings.PricingId, payload pricings.Pricing) *subscriptionPricingUpdatePoller {
	return &subscriptionPricingUpdatePoller{
		client:  client,
		id:      id,
		payload: payload,
	}
}

func (p *subscriptionPricingUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	// the API processes a single pricing update per subscription and returns a 409 until the in-flight one has settled,
	// which the Provider cannot guard against since the previous update may have been made outside of this Terraform run
	resp, err := p.client.Update(ctx, p.id, p.payload)
	if err != nil {
		if response.WasConflict(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return &pollers.PollResult{
			PollInterval: 30 * time.Second,
			Status:       pollers.PollingStatusFailed,
		}, err
	}

	p.Response = resp

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
