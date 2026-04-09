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

func NewResourceProviderFeatureRegistrationPoller(client *features.FeaturesClient, id features.FeatureId, target string) *ResourceProviderFeatureRegistrationPoller {
	return &ResourceProviderFeatureRegistrationPoller{
		client:      client,
		id:          id,
		targetState: target,
	}
}

type ResourceProviderFeatureRegistrationPoller struct {
	client      *features.FeaturesClient
	id          features.FeatureId
	targetState string
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
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 10 * time.Second,
		}, nil
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
