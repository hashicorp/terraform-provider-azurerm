// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileservices"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	storageClients "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

var _ pollers.PollerType = &DataPlaneFileShareAvailabilityPoller{}

const fileServiceResourceNotFoundErrorCode = "ResourceNotFound"

type fileServicesClient interface {
	GetServiceProperties(ctx context.Context, id commonids.StorageAccountId) (fileservices.GetServicePropertiesOperationResponse, error)
}

type DataPlaneFileShareAvailabilityPoller struct {
	client           fileServicesClient
	storageAccountId commonids.StorageAccountId
}

func NewDataPlaneFileShareAvailabilityPoller(client *storageClients.Client, account *storageClients.AccountDetails) (*DataPlaneFileShareAvailabilityPoller, error) {
	return &DataPlaneFileShareAvailabilityPoller{
		client:           client.ResourceManager.FileServices,
		storageAccountId: account.StorageAccountId,
	}, nil
}

func (d *DataPlaneFileShareAvailabilityPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := d.client.GetServiceProperties(ctx, d.storageAccountId)
	serviceIsBeingProvisioned := response.WasNotFound(resp.HttpResponse) &&
		resp.OData != nil &&
		resp.OData.Error != nil &&
		resp.OData.Error.Code != nil &&
		*resp.OData.Error.Code == fileServiceResourceNotFoundErrorCode

	if err != nil {
		if !serviceIsBeingProvisioned {
			return nil, pollers.PollingFailedError{
				Message: err.Error(),
				HttpResponse: &client.Response{
					Response: resp.HttpResponse,
				},
			}
		}
	}
	if serviceIsBeingProvisioned {
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
