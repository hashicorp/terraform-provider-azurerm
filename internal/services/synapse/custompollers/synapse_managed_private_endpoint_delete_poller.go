package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type synapseManagedPrivateEndpointDeletePoller struct {
	client *managedprivateendpoints.ManagedPrivateEndpointsClient
	id     managedprivateendpoints.ManagedPrivateEndpointId
}

var _ pollers.PollerType = &synapseManagedPrivateEndpointDeletePoller{}

func NewSynapseManagedPrivateEndpointDeletePoller(client *managedprivateendpoints.ManagedPrivateEndpointsClient, id managedprivateendpoints.ManagedPrivateEndpointId) pollers.PollerType {
	return &synapseManagedPrivateEndpointDeletePoller{
		client: client,
		id:     id,
	}
}

func (s synapseManagedPrivateEndpointDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := s.client.Get(ctx, s.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				Status: pollers.PollingStatusSucceeded,
			}, nil
		}

		return nil, err
	}

	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
