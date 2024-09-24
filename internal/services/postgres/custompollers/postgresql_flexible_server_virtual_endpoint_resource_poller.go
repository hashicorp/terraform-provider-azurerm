// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &postgresFlexibleServerVirtualEndpointPoller{}

type postgresFlexibleServerVirtualEndpointPoller struct {
	client *virtualendpoints.VirtualEndpointsClient
	id     virtualendpoints.VirtualEndpointId
}

// Workaround due to Azure performing a pseudo-soft delete on virtual endpoints.
//
// - The `DELETE` endpoint does not fully delete the resource, it sets `properties.members` to nil
//
// - Subsequent `GET` operations for the endpoint will always return 200 with empty metadata, so Terraform will hang on `DELETE`
//
// - The only way to currently check for deletion is to check the `properties.members` property
func NewPostgresFlexibleServerVirtualEndpointDeletePoller(client *virtualendpoints.VirtualEndpointsClient, id virtualendpoints.VirtualEndpointId) *postgresFlexibleServerVirtualEndpointPoller {
	return &postgresFlexibleServerVirtualEndpointPoller{
		client: client,
		id:     id,
	}
}

func (p postgresFlexibleServerVirtualEndpointPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	model := resp.Model

	if model != nil && model.Properties != nil {
		if model.Properties.Members != nil {
			return &pollers.PollResult{
				HttpResponse: &client.Response{
					Response: resp.HttpResponse,
				},
				PollInterval: 5 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	return nil, pollers.PollingFailedError{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		Message: fmt.Sprintf("failed to delete %s", p.id),
	}
}
