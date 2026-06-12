// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
)

var _ pollers.PollerType = &DataPlaneTablesAvailabilityPoller{}

type DataPlaneTablesAvailabilityPoller struct {
	client           shim.StorageTableWrapper
	storageAccountId commonids.StorageAccountId
}

func NewDataPlaneTablesAvailabilityPoller(ctx context.Context, client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneTablesAvailabilityPoller, error) {
	tableClient, err := client.TablesDataPlaneClient(ctx, *account, client.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Tables Client: %+v", err)
	}

	return &DataPlaneTablesAvailabilityPoller{
		client:           tableClient,
		storageAccountId: account.StorageAccountId,
	}, nil
}

func (d *DataPlaneTablesAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx)
	var e pollers.PollingDroppedConnectionError
	if errors.As(err, &e) {
		return nil, err
	}
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
