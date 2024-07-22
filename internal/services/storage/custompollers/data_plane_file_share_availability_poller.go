// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

var _ pollers.PollerType = &DataPlaneFileShareAvailabilityPoller{}

type DataPlaneFileShareAvailabilityPoller struct {
	client           *fileservice.FileServiceClient
	storageAccountId commonids.StorageAccountId
}

func NewDataPlaneFileShareAvailabilityPoller(client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneFileShareAvailabilityPoller, error) {
	return &DataPlaneFileShareAvailabilityPoller{
		client:           client.ResourceManager.FileService,
		storageAccountId: account.StorageAccountId,
	}, nil
}

func (d *DataPlaneFileShareAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx, d.storageAccountId)
	if err != nil {
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
