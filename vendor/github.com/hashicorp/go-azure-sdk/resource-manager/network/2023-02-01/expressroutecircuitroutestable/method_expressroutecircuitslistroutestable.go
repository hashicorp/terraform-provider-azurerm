package expressroutecircuitroutestable

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

type ExpressRouteCircuitsListRoutesTableOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ExpressRouteCircuitsListRoutesTable ...
func (c ExpressRouteCircuitRoutesTableClient) ExpressRouteCircuitsListRoutesTable(ctx context.Context, id PeeringRouteTableId) (result ExpressRouteCircuitsListRoutesTableOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
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

// ExpressRouteCircuitsListRoutesTableThenPoll performs ExpressRouteCircuitsListRoutesTable then polls until it's completed
func (c ExpressRouteCircuitRoutesTableClient) ExpressRouteCircuitsListRoutesTableThenPoll(ctx context.Context, id PeeringRouteTableId) error {
	result, err := c.ExpressRouteCircuitsListRoutesTable(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ExpressRouteCircuitsListRoutesTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ExpressRouteCircuitsListRoutesTable: %+v", err)
	}

	return nil
}
