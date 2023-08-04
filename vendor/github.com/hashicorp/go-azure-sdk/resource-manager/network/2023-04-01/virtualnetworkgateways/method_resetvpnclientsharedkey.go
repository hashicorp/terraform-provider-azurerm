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

type ResetVpnClientSharedKeyOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ResetVpnClientSharedKey ...
func (c VirtualNetworkGatewaysClient) ResetVpnClientSharedKey(ctx context.Context, id VirtualNetworkGatewayId) (result ResetVpnClientSharedKeyOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/resetvpnclientsharedkey", id.ID()),
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

// ResetVpnClientSharedKeyThenPoll performs ResetVpnClientSharedKey then polls until it's completed
func (c VirtualNetworkGatewaysClient) ResetVpnClientSharedKeyThenPoll(ctx context.Context, id VirtualNetworkGatewayId) error {
	result, err := c.ResetVpnClientSharedKey(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ResetVpnClientSharedKey: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ResetVpnClientSharedKey: %+v", err)
	}

	return nil
}
