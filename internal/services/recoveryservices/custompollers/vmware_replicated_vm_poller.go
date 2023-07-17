package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &VMWareReplicatedVMPoller{}

func NewVMWareReplicatedVMPoller(client *replicationprotecteditems.ReplicationProtectedItemsClient, id replicationprotecteditems.ReplicationProtectedItemId) *VMWareReplicatedVMPoller {
	return &VMWareReplicatedVMPoller{
		client: client,
		id:     id,
	}
}

type VMWareReplicatedVMPoller struct {
	client *replicationprotecteditems.ReplicationProtectedItemsClient
	id     replicationprotecteditems.ReplicationProtectedItemId
}

func (p VMWareReplicatedVMPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	protectionState := ""
	if model := resp.Model; model != nil && model.Properties != nil && resp.Model.Properties.ProtectionState != nil {
		protectionState = *model.Properties.ProtectionState
	}

	if strings.EqualFold(protectionState, "Protected") || strings.EqualFold(protectionState, "normal") {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 1 * time.Minute,
		}, nil
	}

	// Processing
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 1 * time.Minute,
	}, nil
}
