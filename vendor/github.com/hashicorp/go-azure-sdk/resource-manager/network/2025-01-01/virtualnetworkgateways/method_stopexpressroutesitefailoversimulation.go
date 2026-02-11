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

type StopExpressRouteSiteFailoverSimulationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

// StopExpressRouteSiteFailoverSimulation ...
func (c VirtualNetworkGatewaysClient) StopExpressRouteSiteFailoverSimulation(ctx context.Context, id VirtualNetworkGatewayId, input ExpressRouteFailoverStopApiParameters) (result StopExpressRouteSiteFailoverSimulationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/stopSiteFailoverTest", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// StopExpressRouteSiteFailoverSimulationThenPoll performs StopExpressRouteSiteFailoverSimulation then polls until it's completed
func (c VirtualNetworkGatewaysClient) StopExpressRouteSiteFailoverSimulationThenPoll(ctx context.Context, id VirtualNetworkGatewayId, input ExpressRouteFailoverStopApiParameters) error {
	result, err := c.StopExpressRouteSiteFailoverSimulation(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing StopExpressRouteSiteFailoverSimulation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StopExpressRouteSiteFailoverSimulation: %+v", err)
	}

	return nil
}
