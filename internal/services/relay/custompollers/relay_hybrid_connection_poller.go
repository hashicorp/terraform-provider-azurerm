package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &relayHybridConnectionPoller{}

type relayHybridConnectionPoller struct {
	client *hybridconnections.HybridConnectionsClient
	id     hybridconnections.HybridConnectionId
}

var (
	relayHybridConnectionPollingSuccess = pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	relayHybridConnectionPollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func DeleteRelayHybridConnectionPoller(client *hybridconnections.HybridConnectionsClient, id hybridconnections.HybridConnectionId) *relayHybridConnectionPoller {
	return &relayHybridConnectionPoller{
		client: client,
		id:     id,
	}
}

func (p relayHybridConnectionPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &relayHybridConnectionPollingSuccess, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &relayHybridConnectionPollingInProgress, nil
}
