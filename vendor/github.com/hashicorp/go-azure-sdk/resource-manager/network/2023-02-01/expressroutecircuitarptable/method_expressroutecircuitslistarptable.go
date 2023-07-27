package expressroutecircuitarptable

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

type ExpressRouteCircuitsListArpTableOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ExpressRouteCircuitsListArpTable ...
func (c ExpressRouteCircuitArpTableClient) ExpressRouteCircuitsListArpTable(ctx context.Context, id ArpTableId) (result ExpressRouteCircuitsListArpTableOperationResponse, err error) {
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

// ExpressRouteCircuitsListArpTableThenPoll performs ExpressRouteCircuitsListArpTable then polls until it's completed
func (c ExpressRouteCircuitArpTableClient) ExpressRouteCircuitsListArpTableThenPoll(ctx context.Context, id ArpTableId) error {
	result, err := c.ExpressRouteCircuitsListArpTable(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ExpressRouteCircuitsListArpTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ExpressRouteCircuitsListArpTable: %+v", err)
	}

	return nil
}
