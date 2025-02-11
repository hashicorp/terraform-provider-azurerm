// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

var _ pollers.PollerType = &DataPlaneBlobContainersAvailabilityPoller{}

type DataPlaneBlobContainersAvailabilityPoller struct {
	client      *accounts.Client
	accountName string
}

func NewDataPlaneBlobContainersAvailabilityPoller(ctx context.Context, client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneBlobContainersAvailabilityPoller, error) {
	dataPlaneClient, err := client.AccountsDataPlaneClient(ctx, *account, client.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, err
	}
	return &DataPlaneBlobContainersAvailabilityPoller{
		client:      dataPlaneClient,
		accountName: account.StorageAccountId.StorageAccountName,
	}, nil
}

func (d *DataPlaneBlobContainersAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx, d.accountName)
	if err != nil {
		if resp.HttpResponse == nil {
			return nil, pollers.PollingDroppedConnectionError{
				Message: err.Error(),
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return nil, pollers.PollingFailedError{
				Message: err.Error(),
				HttpResponse: &client.Response{
					Response: resp.HttpResponse,
				},
			}
		}
	}
	if response.WasNotFound(resp.HttpResponse) {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 10 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
