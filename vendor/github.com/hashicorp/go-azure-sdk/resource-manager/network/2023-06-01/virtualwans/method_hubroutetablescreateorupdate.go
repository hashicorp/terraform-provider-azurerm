package virtualwans

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

type HubRouteTablesCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// HubRouteTablesCreateOrUpdate ...
func (c VirtualWANsClient) HubRouteTablesCreateOrUpdate(ctx context.Context, id HubRouteTableId, input HubRouteTable) (result HubRouteTablesCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
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

// HubRouteTablesCreateOrUpdateThenPoll performs HubRouteTablesCreateOrUpdate then polls until it's completed
func (c VirtualWANsClient) HubRouteTablesCreateOrUpdateThenPoll(ctx context.Context, id HubRouteTableId, input HubRouteTable) error {
	result, err := c.HubRouteTablesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing HubRouteTablesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after HubRouteTablesCreateOrUpdate: %+v", err)
	}

	return nil
}
