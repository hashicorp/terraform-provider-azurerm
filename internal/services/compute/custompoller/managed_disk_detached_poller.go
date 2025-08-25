package custompoller

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type managedDiskDetachedPoller struct {
	client *disks.DisksClient
	id     commonids.ManagedDiskId
	vmId   virtualmachines.VirtualMachineId
}

var _ pollers.PollerType = &managedDiskDetachedPoller{}

var (
	pollingSuccess = pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}

	pollingFailed = pollers.PollResult{
		Status: pollers.PollingStatusFailed,
	}

	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}

	pollingUnknown = pollers.PollResult{
		Status:       pollers.PollingStatusUnknown,
		PollInterval: 10 * time.Second,
	}
)

func NewManagedDiskDetachedPoller(client *disks.DisksClient, id commonids.ManagedDiskId, vmId virtualmachines.VirtualMachineId) *managedDiskDetachedPoller {
	return &managedDiskDetachedPoller{
		client: client,
		id:     id,
		vmId:   vmId,
	}
}

func (m managedDiskDetachedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := m.client.Get(ctx, m.id)
	if err != nil {
		return &pollingFailed, fmt.Errorf("checking Managed Disk (%s) for detach operation during deletion of %s: %+v", m.id, m.vmId, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil {
		if pointer.From(resp.Model.Properties.DiskState) == disks.DiskStateAttached {
			return &pollingInProgress, nil
		}
	} else {
		return &pollingUnknown, nil
	}

	return &pollingSuccess, nil
}
