// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var (
	_ pollers.PollerType = &threatIntelligenceIndicatorPoller{}
)

const consistentRequestCount = 10

type threatIntelligenceIndicatorPoller struct {
	client              *threatintelligence.ThreatIntelligenceClient
	id                  threatintelligence.IndicatorId
	successfulPollCount int
}

func NewThreatIntelligenceIndicatorPoller(client *threatintelligence.ThreatIntelligenceClient, id threatintelligence.IndicatorId) *threatIntelligenceIndicatorPoller {
	return &threatIntelligenceIndicatorPoller{
		client: client,
		id:     id,
	}
}

func (p *threatIntelligenceIndicatorPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.IndicatorGet(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 5 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if p.successfulPollCount < consistentRequestCount {
		p.successfulPollCount++
		return &pollers.PollResult{
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
