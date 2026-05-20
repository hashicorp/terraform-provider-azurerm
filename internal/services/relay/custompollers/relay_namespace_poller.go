package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &relayNamespacePoller{}

type relayNamespacePoller struct {
	client *namespaces.NamespacesClient
	id     namespaces.NamespaceId
}

var (
	relayNamespacePollingSuccess = pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	relayNamespacePollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func DeleteRelayNamespacePoller(client *namespaces.NamespacesClient, id namespaces.NamespaceId) *relayNamespacePoller {
	return &relayNamespacePoller{
		client: client,
		id:     id,
	}
}

func (p relayNamespacePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &relayNamespacePollingSuccess, nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &relayNamespacePollingInProgress, nil
}
