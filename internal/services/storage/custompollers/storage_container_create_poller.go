// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &storageContainerCreatePoller{}

type storageContainerCreatePoller struct {
	client          *blobcontainers.BlobContainersClient
	id              commonids.StorageContainerId
	initialResponse *http.Response
}

func NewStorageContainerCreatePoller(client *blobcontainers.BlobContainersClient, id commonids.StorageContainerId, initialResponse *http.Response) *storageContainerCreatePoller {
	return &storageContainerCreatePoller{
		client:          client,
		id:              id,
		initialResponse: initialResponse,
	}
}

func (p *storageContainerCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if p.initialResponse != nil {
		resp := p.initialResponse
		p.initialResponse = nil
		return &pollers.PollResult{
			HttpResponse: &client.Response{Response: resp},
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	result, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return &pollers.PollResult{
				HttpResponse: &client.Response{Response: result.HttpResponse},
				PollInterval: 5 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}
		return &pollers.PollResult{
			HttpResponse: &client.Response{Response: result.HttpResponse},
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusFailed,
		}, err
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{Response: result.HttpResponse},
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
