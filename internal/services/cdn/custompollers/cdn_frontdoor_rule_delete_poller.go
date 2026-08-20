// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorRuleDeletePoller{}

type frontDoorRuleDeletePoller struct {
	client *rules.RulesClient
	id     rules.RuleId
}

func NewFrontDoorRuleDeletePoller(client *rules.RulesClient, id rules.RuleId) pollers.PollerType {
	return &frontDoorRuleDeletePoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorRuleDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusSucceeded,
			}, nil
		}

		return nil, fmt.Errorf("checking deletion of %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
