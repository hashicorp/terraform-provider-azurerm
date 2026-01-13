package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type dataFactoryIntegrationRuntimeStatusPoller struct {
	client *integrationruntimes.IntegrationRuntimesClient
	id     integrationruntimes.IntegrationRuntimeId
}

var _ pollers.PollerType = &dataFactoryIntegrationRuntimeStatusPoller{}

func (p dataFactoryIntegrationRuntimeStatusPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result := &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}

	resp, err := p.client.GetStatus(ctx, p.id)
	if err != nil {
		result.Status = pollers.PollingStatusFailed
		return result, pollers.PollingFailedError{Message: fmt.Sprintf("retrieving status for %s: %+v", p.id, err)}
	}

	if resp.Model != nil && pointer.From(resp.Model.Properties.IntegrationRuntimeStatus().State) == integrationruntimes.IntegrationRuntimeStateOnline {
		result.Status = pollers.PollingStatusSucceeded
		return result, nil
	}

	return result, nil
}

func NewDataFactoryIntegrationRuntimeStatusPoller(client *integrationruntimes.IntegrationRuntimesClient, id integrationruntimes.IntegrationRuntimeId) pollers.PollerType {
	return &dataFactoryIntegrationRuntimeStatusPoller{
		client: client,
		id:     id,
	}
}
