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

	if *resp.Model.Properties.ProvisioningState == "Succeeded" {
		return &pollers.PollResult{
			Status: pollers.PollingStatusSucceeded,
		}, nil
	}

	// The API spec doesn't document the provisioning states, so unless it's `Succeeded` we'll continue polling
	// ideally we'd exit early on any terminal states that are not `Succeeded` but they are unknown.
	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
