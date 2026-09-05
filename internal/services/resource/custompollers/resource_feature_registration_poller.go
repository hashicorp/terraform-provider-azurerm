// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &ResourceProviderFeatureRegistrationPoller{}

// NewResourceProviderFeatureRegistrationPoller - Polls for an expected registration state, accounting for
// AFEC API replica inconsistency by requiring n continuous occurrences of the target state, in the same
// way as ResourceProviderRegistrationPoller.
func NewResourceProviderFeatureRegistrationPoller(client *features.FeaturesClient, id features.FeatureId, target string) *ResourceProviderFeatureRegistrationPoller {
	return &ResourceProviderFeatureRegistrationPoller{
		client:                     client,
		id:                         id,
		targetState:                target,
		continuousTargetOccurrence: 10,
	}
}

type ResourceProviderFeatureRegistrationPoller struct {
	client                     *features.FeaturesClient
	id                         features.FeatureId
	targetState                string
	continuousTargetOccurrence int
	currentOccurrenceCount     int
}

func (p *ResourceProviderFeatureRegistrationPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.State == nil {
		return nil, fmt.Errorf("retrieving %s: unable to determine registration state", p.id)
	}

	if strings.EqualFold(*resp.Model.Properties.State, p.targetState) {
		if p.currentOccurrenceCount >= p.continuousTargetOccurrence {
			return &pollers.PollResult{
				Status:       pollers.PollingStatusSucceeded,
				PollInterval: 10 * time.Second,
			}, nil
		}
		p.currentOccurrenceCount++
	} else {
		p.currentOccurrenceCount = 0
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
