package custompollers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &networkManagerIPAMPoolPoller{}

type networkManagerIPAMPoolPoller struct {
	client *ipampools.IPamPoolsClient
	id     ipampools.IPamPoolId
}

func NewNetworkManagerIPAMPoolPoller(client *ipampools.IPamPoolsClient, id ipampools.IPamPoolId) *networkManagerIPAMPoolPoller {
	return &networkManagerIPAMPoolPoller{
		client: client,
		id:     id,
	}
}

func (p networkManagerIPAMPoolPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingSuccess, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollingInProgress, nil
}
