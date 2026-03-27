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

type StartExpressRouteSiteFailoverSimulationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

type StartExpressRouteSiteFailoverSimulationOperationOptions struct {
	PeeringLocation *string
}

func DefaultStartExpressRouteSiteFailoverSimulationOperationOptions() StartExpressRouteSiteFailoverSimulationOperationOptions {
	return StartExpressRouteSiteFailoverSimulationOperationOptions{}
}

func (o StartExpressRouteSiteFailoverSimulationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StartExpressRouteSiteFailoverSimulationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StartExpressRouteSiteFailoverSimulationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PeeringLocation != nil {
		out.Append("peeringLocation", fmt.Sprintf("%v", *o.PeeringLocation))
	}
	return &out
}

// StartExpressRouteSiteFailoverSimulation ...
func (c VirtualNetworkGatewaysClient) StartExpressRouteSiteFailoverSimulation(ctx context.Context, id VirtualNetworkGatewayId, options StartExpressRouteSiteFailoverSimulationOperationOptions) (result StartExpressRouteSiteFailoverSimulationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/startSiteFailoverTest", id.ID()),
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

// StartExpressRouteSiteFailoverSimulationThenPoll performs StartExpressRouteSiteFailoverSimulation then polls until it's completed
func (c VirtualNetworkGatewaysClient) StartExpressRouteSiteFailoverSimulationThenPoll(ctx context.Context, id VirtualNetworkGatewayId, options StartExpressRouteSiteFailoverSimulationOperationOptions) error {
	result, err := c.StartExpressRouteSiteFailoverSimulation(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing StartExpressRouteSiteFailoverSimulation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StartExpressRouteSiteFailoverSimulation: %+v", err)
	}

	return nil
}
