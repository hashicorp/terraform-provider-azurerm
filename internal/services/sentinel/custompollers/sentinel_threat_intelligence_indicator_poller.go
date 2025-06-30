// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// The GET request after creation may be returned with HTTP 404 in a period.
// Tracked by https://github.com/Azure/azure-rest-api-specs/issues/35551
var (
	_ pollers.PollerType = &threatIntelligenceIndicatorPooler{}
	_ pollers.PollerType = &threatIntelligenceIndicatorUpdatePooler{}
)

const concistentRequestNumber = 10

type threatIntelligenceIndicatorPooler struct {
	client       azuresdkhacks.ThreatIntelligenceIndicatorClient
	id           parse.ThreatIntelligenceIndicatorId
	succeededCnt int
}

type threatIntelligenceIndicatorUpdatePooler struct {
	client         azuresdkhacks.ThreatIntelligenceIndicatorClient
	id             parse.ThreatIntelligenceIndicatorId
	succeededCnt   int
	lastUpdateTime string
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}

	pollingFailed = pollers.PollResult{
		Status: pollers.PollingStatusFailed,
	}

	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewThreatIntelligenceIndicatorPoller(client azuresdkhacks.ThreatIntelligenceIndicatorClient, id parse.ThreatIntelligenceIndicatorId) *threatIntelligenceIndicatorPooler {
	return &threatIntelligenceIndicatorPooler{
		client:       client,
		id:           id,
		succeededCnt: 0,
	}
}

func (p *threatIntelligenceIndicatorPooler) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id.ResourceGroup, p.id.WorkspaceName, p.id.IndicatorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return &pollingInProgress, nil
		}
		return nil, fmt.Errorf("retrieving %s, %+v", p.id, err)
	}

	if p.succeededCnt < concistentRequestNumber {
		p.succeededCnt++
		return &pollingInProgress, nil
	}

	return &pollingSuccess, nil
}

func NewThreatIntelligenceIndicatorUpdatePoller(client azuresdkhacks.ThreatIntelligenceIndicatorClient, id parse.ThreatIntelligenceIndicatorId, lastUpdateTime string) *threatIntelligenceIndicatorUpdatePooler {
	return &threatIntelligenceIndicatorUpdatePooler{
		client:         client,
		id:             id,
		succeededCnt:   0,
		lastUpdateTime: lastUpdateTime,
	}
}

func (p *threatIntelligenceIndicatorUpdatePooler) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id.ResourceGroup, p.id.WorkspaceName, p.id.IndicatorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return &pollingInProgress, nil
		}
		return &pollingFailed, fmt.Errorf("retrieving %s, %+v", p.id, err)
	}

	model, ok := resp.Value.AsThreatIntelligenceIndicatorModel()
	if !ok {
		return &pollingFailed, fmt.Errorf("retrieving %s: type mismatch", p.id)
	}

	if model.LastUpdatedTimeUtc != nil && model.LastUpdatedTimeUtc != &p.lastUpdateTime {
		return &pollingSuccess, nil
	}

	return &pollingInProgress, nil
}
