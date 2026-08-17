package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &ResourceProviderRegistrationPoller{}

// NewResourceProviderRegistrationPollerDefault - Polls for an expected registration state, accounting for API inconsistency by requiring 10 continuous occurrences
func NewResourceProviderRegistrationPollerDefault(client *providers.ProvidersClient, id providers.SubscriptionProviderId, target string) *ResourceProviderRegistrationPoller {
	return NewResourceProviderRegistrationPoller(client, id, target, 10)
}

// NewResourceProviderRegistrationPoller - Polls for an expected registration state, accounting for API inconsistency by requiring n continuous occurrences
func NewResourceProviderRegistrationPoller(client *providers.ProvidersClient, id providers.SubscriptionProviderId, target string, continuousTargetOccurrence int) *ResourceProviderRegistrationPoller {
	return &ResourceProviderRegistrationPoller{
		client:                     client,
		id:                         id,
		targetState:                target,
		continuousTargetOccurrence: continuousTargetOccurrence,
	}
}

type ResourceProviderRegistrationPoller struct {
	client                     *providers.ProvidersClient
	id                         providers.SubscriptionProviderId
	targetState                string
	continuousTargetOccurrence int
	currentOccurrenceCount     int
}

func (p *ResourceProviderRegistrationPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id, providers.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	registrationState := ""
	if model := resp.Model; model != nil && model.RegistrationState != nil {
		registrationState = *model.RegistrationState
	}

	if strings.EqualFold(registrationState, p.targetState) {
		if p.continuousTargetOccurrence == p.currentOccurrenceCount {
			return &pollers.PollResult{
				Status:       pollers.PollingStatusSucceeded,
				PollInterval: 10 * time.Second,
			}, nil
		}
		p.currentOccurrenceCount += 1
	} else {
		p.currentOccurrenceCount = 0
	}

	// Processing
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
