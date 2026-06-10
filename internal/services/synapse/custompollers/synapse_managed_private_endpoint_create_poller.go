package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type synapseManagedPrivateEndpointCreatePoller struct {
	client *managedprivateendpoints.ManagedPrivateEndpointsClient
	id     managedprivateendpoints.ManagedPrivateEndpointId
}

var _ pollers.PollerType = &synapseManagedPrivateEndpointCreatePoller{}

func NewSynapseManagedPrivateEndpointCreatePoller(client *managedprivateendpoints.ManagedPrivateEndpointsClient, id managedprivateendpoints.ManagedPrivateEndpointId) pollers.PollerType {
	return &synapseManagedPrivateEndpointCreatePoller{
		client: client,
		id:     id,
	}
}

func (s synapseManagedPrivateEndpointCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := s.client.Get(ctx, s.id)
	if err != nil {
		return nil, err
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.ProvisioningState == nil {
		return nil, fmt.Errorf("checking `provisioningState` for %s", s.id)
	}

	switch *resp.Model.Properties.ProvisioningState {
	case string(pollers.PollingStatusSucceeded):
		return &pollers.PollResult{
			Status: pollers.PollingStatusSucceeded,
		}, nil
	case string(pollers.PollingStatusFailed):
		return nil, fmt.Errorf("provisioningState was `%s`", pollers.PollingStatusFailed)
	case string(pollers.PollingStatusCancelled):
		return nil, fmt.Errorf("provisioningState was `%s`", pollers.PollingStatusCancelled)
	}

	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
