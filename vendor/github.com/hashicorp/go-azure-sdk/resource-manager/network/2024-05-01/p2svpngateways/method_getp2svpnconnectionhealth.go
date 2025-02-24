package p2svpngateways

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

type GetP2sVpnConnectionHealthOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *P2SVpnGateway
}

// GetP2sVpnConnectionHealth ...
func (c P2sVpnGatewaysClient) GetP2sVpnConnectionHealth(ctx context.Context, id commonids.VirtualWANP2SVPNGatewayId) (result GetP2sVpnConnectionHealthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/getP2sVpnConnectionHealth", id.ID()),
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

// GetP2sVpnConnectionHealthThenPoll performs GetP2sVpnConnectionHealth then polls until it's completed
func (c P2sVpnGatewaysClient) GetP2sVpnConnectionHealthThenPoll(ctx context.Context, id commonids.VirtualWANP2SVPNGatewayId) error {
	result, err := c.GetP2sVpnConnectionHealth(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GetP2sVpnConnectionHealth: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetP2sVpnConnectionHealth: %+v", err)
	}

	return nil
}
