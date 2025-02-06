package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubVirtualNetworkConnectionsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// HubVirtualNetworkConnectionsDelete ...
func (c VirtualWANsClient) HubVirtualNetworkConnectionsDelete(ctx context.Context, id HubVirtualNetworkConnectionId) (result HubVirtualNetworkConnectionsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// HubVirtualNetworkConnectionsDeleteThenPoll performs HubVirtualNetworkConnectionsDelete then polls until it's completed
func (c VirtualWANsClient) HubVirtualNetworkConnectionsDeleteThenPoll(ctx context.Context, id HubVirtualNetworkConnectionId) error {
	result, err := c.HubVirtualNetworkConnectionsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing HubVirtualNetworkConnectionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after HubVirtualNetworkConnectionsDelete: %+v", err)
	}

	return nil
}
