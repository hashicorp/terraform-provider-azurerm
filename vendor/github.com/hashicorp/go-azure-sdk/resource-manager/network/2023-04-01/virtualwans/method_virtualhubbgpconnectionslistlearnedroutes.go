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

type VirtualHubBgpConnectionsListLearnedRoutesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// VirtualHubBgpConnectionsListLearnedRoutes ...
func (c VirtualWANsClient) VirtualHubBgpConnectionsListLearnedRoutes(ctx context.Context, id commonids.VirtualHubBGPConnectionId) (result VirtualHubBgpConnectionsListLearnedRoutesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/learnedRoutes", id.ID()),
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

// VirtualHubBgpConnectionsListLearnedRoutesThenPoll performs VirtualHubBgpConnectionsListLearnedRoutes then polls until it's completed
func (c VirtualWANsClient) VirtualHubBgpConnectionsListLearnedRoutesThenPoll(ctx context.Context, id commonids.VirtualHubBGPConnectionId) error {
	result, err := c.VirtualHubBgpConnectionsListLearnedRoutes(ctx, id)
	if err != nil {
		return fmt.Errorf("performing VirtualHubBgpConnectionsListLearnedRoutes: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after VirtualHubBgpConnectionsListLearnedRoutes: %+v", err)
	}

	return nil
}
