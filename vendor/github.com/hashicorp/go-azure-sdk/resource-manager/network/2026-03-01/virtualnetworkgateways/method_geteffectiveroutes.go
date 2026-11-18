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

type GetEffectiveRoutesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GatewayEffectiveRouteListResult
}

// GetEffectiveRoutes ...
func (c VirtualNetworkGatewaysClient) GetEffectiveRoutes(ctx context.Context, id VirtualNetworkGatewayId) (result GetEffectiveRoutesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/getEffectiveRoutes", id.ID()),
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

// GetEffectiveRoutesThenPoll performs GetEffectiveRoutes then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetEffectiveRoutesThenPoll(ctx context.Context, id VirtualNetworkGatewayId) error {
	return c.GetEffectiveRoutesCallbackThenPoll(ctx, id, nil)
}

// GetEffectiveRoutesCallbackThenPoll performs GetEffectiveRoutes, runs the optional callback function, then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetEffectiveRoutesCallbackThenPoll(ctx context.Context, id VirtualNetworkGatewayId, callback func() error) error {
	result, err := c.GetEffectiveRoutes(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GetEffectiveRoutes: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetEffectiveRoutes: %+v", err)
	}

	return nil
}
