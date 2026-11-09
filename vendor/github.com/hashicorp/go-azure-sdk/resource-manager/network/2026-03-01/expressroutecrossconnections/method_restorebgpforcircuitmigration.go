package expressroutecrossconnections

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

type RestoreBgpForCircuitMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MigrateExpressRouteCircuitHealthCheckResponse
}

// RestoreBgpForCircuitMigration ...
func (c ExpressRouteCrossConnectionsClient) RestoreBgpForCircuitMigration(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitRequest) (result RestoreBgpForCircuitMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/restoreBgpForCircuitMigration", id.ID()),
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

// RestoreBgpForCircuitMigrationThenPoll performs RestoreBgpForCircuitMigration then polls until it's completed
func (c ExpressRouteCrossConnectionsClient) RestoreBgpForCircuitMigrationThenPoll(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitRequest) error {
	return c.RestoreBgpForCircuitMigrationCallbackThenPoll(ctx, id, input, nil)
}

// RestoreBgpForCircuitMigrationCallbackThenPoll performs RestoreBgpForCircuitMigration, runs the optional callback function, then polls until it's completed
func (c ExpressRouteCrossConnectionsClient) RestoreBgpForCircuitMigrationCallbackThenPoll(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitRequest, callback func() error) error {
	result, err := c.RestoreBgpForCircuitMigration(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RestoreBgpForCircuitMigration: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RestoreBgpForCircuitMigration: %+v", err)
	}

	return nil
}
