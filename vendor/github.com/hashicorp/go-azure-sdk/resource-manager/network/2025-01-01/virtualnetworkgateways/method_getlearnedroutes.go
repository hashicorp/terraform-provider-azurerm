package virtualnetworkgateways

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

type GetLearnedRoutesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GatewayRouteListResult
}

// GetLearnedRoutes ...
func (c VirtualNetworkGatewaysClient) GetLearnedRoutes(ctx context.Context, id VirtualNetworkGatewayId) (result GetLearnedRoutesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/getLearnedRoutes", id.ID()),
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

// GetLearnedRoutesThenPoll performs GetLearnedRoutes then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetLearnedRoutesThenPoll(ctx context.Context, id VirtualNetworkGatewayId) error {
	result, err := c.GetLearnedRoutes(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GetLearnedRoutes: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetLearnedRoutes: %+v", err)
	}

	return nil
}
