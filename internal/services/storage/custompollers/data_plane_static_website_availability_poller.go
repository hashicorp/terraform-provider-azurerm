// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

var _ pollers.PollerType = &DataPlaneStaticWebsiteAvailabilityPoller{}

type DataPlaneStaticWebsiteAvailabilityPoller struct {
	client           *accounts.Client
	storageAccountId commonids.StorageAccountId
}

func NewDataPlaneStaticWebsiteAvailabilityPoller(ctx context.Context, client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneStaticWebsiteAvailabilityPoller, error) {
	accountsClient, err := client.AccountsDataPlaneClient(ctx, *account, client.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Accounts Data Plane Client: %+v", err)
	}

	return &DataPlaneStaticWebsiteAvailabilityPoller{
		client:           accountsClient,
		storageAccountId: account.StorageAccountId,
	}, nil
}

func (d *DataPlaneStaticWebsiteAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx, d.storageAccountId.StorageAccountName)
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
