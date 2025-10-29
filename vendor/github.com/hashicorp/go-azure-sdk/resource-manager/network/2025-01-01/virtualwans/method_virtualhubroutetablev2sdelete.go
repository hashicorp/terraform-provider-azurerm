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

type VirtualHubRouteTableV2sDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// VirtualHubRouteTableV2sDelete ...
func (c VirtualWANsClient) VirtualHubRouteTableV2sDelete(ctx context.Context, id VirtualHubRouteTableId) (result VirtualHubRouteTableV2sDeleteOperationResponse, err error) {
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

// VirtualHubRouteTableV2sDeleteThenPoll performs VirtualHubRouteTableV2sDelete then polls until it's completed
func (c VirtualWANsClient) VirtualHubRouteTableV2sDeleteThenPoll(ctx context.Context, id VirtualHubRouteTableId) error {
	result, err := c.VirtualHubRouteTableV2sDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing VirtualHubRouteTableV2sDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after VirtualHubRouteTableV2sDelete: %+v", err)
	}

	return nil
}
