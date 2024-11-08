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

type GetAdvertisedRoutesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GatewayRouteListResult
}

type GetAdvertisedRoutesOperationOptions struct {
	Peer *string
}

func DefaultGetAdvertisedRoutesOperationOptions() GetAdvertisedRoutesOperationOptions {
	return GetAdvertisedRoutesOperationOptions{}
}

func (o GetAdvertisedRoutesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetAdvertisedRoutesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetAdvertisedRoutesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Peer != nil {
		out.Append("peer", fmt.Sprintf("%v", *o.Peer))
	}
	return &out
}

// GetAdvertisedRoutes ...
func (c VirtualNetworkGatewaysClient) GetAdvertisedRoutes(ctx context.Context, id VirtualNetworkGatewayId, options GetAdvertisedRoutesOperationOptions) (result GetAdvertisedRoutesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getAdvertisedRoutes", id.ID()),
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

// GetAdvertisedRoutesThenPoll performs GetAdvertisedRoutes then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetAdvertisedRoutesThenPoll(ctx context.Context, id VirtualNetworkGatewayId, options GetAdvertisedRoutesOperationOptions) error {
	result, err := c.GetAdvertisedRoutes(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetAdvertisedRoutes: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetAdvertisedRoutes: %+v", err)
	}

	return nil
}
