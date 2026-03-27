package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type resourceGroupCreatePoller struct {
	client *resourcegroups.ResourceGroupsClient
	id     commonids.ResourceGroupId
}

var _ pollers.PollerType = &resourceGroupCreatePoller{}

var (
	successCount   = 3 // emulates ContinuousTargetOccurrence
	pollingSuccess = &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
	pollingFailed = &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusFailed,
	}
)

func NewResourceGroupCreatePoller(client *resourcegroups.ResourceGroupsClient, id commonids.ResourceGroupId) *resourceGroupCreatePoller {
	return &resourceGroupCreatePoller{
		client: client,
		id:     id,
	}
}

func (p resourceGroupCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	rg, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(rg.HttpResponse) {
			successCount = 3
			return pollingInProgress, nil
		}
		return pollingFailed, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if successCount > 1 {
		successCount--
		return pollingInProgress, nil
	}

	return pollingSuccess, nil
}
