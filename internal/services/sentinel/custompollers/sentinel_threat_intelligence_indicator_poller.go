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
	_ pollers.PollerType = &threatIntelligenceIndicatorUpdatePoller{}
)

const consistentRequestCount = 10

type threatIntelligenceIndicatorPoller struct {
	client              *threatintelligence.ThreatIntelligenceClient
	id                  threatintelligence.IndicatorId
	successfulPollCount int
}

type threatIntelligenceIndicatorUpdatePoller struct {
	client              *threatintelligence.ThreatIntelligenceClient
	id                  threatintelligence.IndicatorId
	previousEtag        string
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

func NewThreatIntelligenceIndicatorUpdatePoller(client *threatintelligence.ThreatIntelligenceClient, id threatintelligence.IndicatorId, previousEtag string) *threatIntelligenceIndicatorUpdatePoller {
	return &threatIntelligenceIndicatorUpdatePoller{
		client:       client,
		id:           id,
		previousEtag: previousEtag,
	}
}

func (p *threatIntelligenceIndicatorUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
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

	model, ok := resp.Model.(threatintelligence.ThreatIntelligenceIndicatorModel)
	if !ok {
		return nil, fmt.Errorf("retrieving %s: type mismatch, got %T", p.id, resp.Model)
	}

	if !threatIntelligenceIndicatorEtagUpdated(model.Etag, p.previousEtag) {
		p.successfulPollCount = 0
		return &pollers.PollResult{
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	p.successfulPollCount++
	status := pollers.PollingStatusInProgress
	if p.successfulPollCount >= consistentRequestCount {
		status = pollers.PollingStatusSucceeded
	}

	return &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       status,
	}, nil
}

func threatIntelligenceIndicatorEtagUpdated(actual *string, previous string) bool {
	return actual != nil && (previous == "" || *actual != previous)
}
