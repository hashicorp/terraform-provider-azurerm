package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type P2sVpnGatewaysDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// P2sVpnGatewaysDelete ...
func (c VirtualWANsClient) P2sVpnGatewaysDelete(ctx context.Context, id commonids.VirtualWANP2SVPNGatewayId) (result P2sVpnGatewaysDeleteOperationResponse, err error) {
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

// P2sVpnGatewaysDeleteThenPoll performs P2sVpnGatewaysDelete then polls until it's completed
func (c VirtualWANsClient) P2sVpnGatewaysDeleteThenPoll(ctx context.Context, id commonids.VirtualWANP2SVPNGatewayId) error {
	result, err := c.P2sVpnGatewaysDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing P2sVpnGatewaysDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after P2sVpnGatewaysDelete: %+v", err)
	}

	return nil
}
