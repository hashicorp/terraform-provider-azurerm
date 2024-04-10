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

type GetBgpPeerStatusOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *BgpPeerStatusListResult
}

type GetBgpPeerStatusOperationOptions struct {
	Peer *string
}

func DefaultGetBgpPeerStatusOperationOptions() GetBgpPeerStatusOperationOptions {
	return GetBgpPeerStatusOperationOptions{}
}

func (o GetBgpPeerStatusOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetBgpPeerStatusOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetBgpPeerStatusOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Peer != nil {
		out.Append("peer", fmt.Sprintf("%v", *o.Peer))
	}
	return &out
}

// GetBgpPeerStatus ...
func (c VirtualNetworkGatewaysClient) GetBgpPeerStatus(ctx context.Context, id VirtualNetworkGatewayId, options GetBgpPeerStatusOperationOptions) (result GetBgpPeerStatusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/getBgpPeerStatus", id.ID()),
		OptionsObject: options,
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

// GetBgpPeerStatusThenPoll performs GetBgpPeerStatus then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetBgpPeerStatusThenPoll(ctx context.Context, id VirtualNetworkGatewayId, options GetBgpPeerStatusOperationOptions) error {
	result, err := c.GetBgpPeerStatus(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetBgpPeerStatus: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetBgpPeerStatus: %+v", err)
	}

	return nil
}
