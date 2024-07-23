// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
)

var _ pollers.PollerType = &DataPlaneQueuesAvailabilityPoller{}

type DataPlaneQueuesAvailabilityPoller struct {
	client           shim.StorageQueuesWrapper
	storageAccountId commonids.StorageAccountId
}

func NewDataPlaneQueuesAvailabilityPoller(ctx context.Context, client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneQueuesAvailabilityPoller, error) {
	queueClient, err := client.QueuesDataPlaneClient(ctx, *account, client.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Queues Client: %+v", err)
	}

	return &DataPlaneQueuesAvailabilityPoller{
		client:           queueClient,
		storageAccountId: account.StorageAccountId,
	}, nil
}

func (d *DataPlaneQueuesAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx)
	if err != nil {
		return nil, pollers.PollingFailedError{
			Message:      err.Error(),
			HttpResponse: nil,
		}
	}
	if resp == nil {
		return &pollers.PollResult{
			HttpResponse: nil,
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
