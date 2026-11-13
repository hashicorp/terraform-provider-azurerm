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

type GetCircuitMigrationInfoOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MigrateExpressRouteCircuitHealthCheckResponse
}

// GetCircuitMigrationInfo ...
func (c ExpressRouteCrossConnectionsClient) GetCircuitMigrationInfo(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitValidateAndHealthCheckRequest) (result GetCircuitMigrationInfoOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/getCircuitMigrationInfo", id.ID()),
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

// GetCircuitMigrationInfoThenPoll performs GetCircuitMigrationInfo then polls until it's completed
func (c ExpressRouteCrossConnectionsClient) GetCircuitMigrationInfoThenPoll(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitValidateAndHealthCheckRequest) error {
	return c.GetCircuitMigrationInfoCallbackThenPoll(ctx, id, input, nil)
}

// GetCircuitMigrationInfoCallbackThenPoll performs GetCircuitMigrationInfo, runs the optional callback function, then polls until it's completed
func (c ExpressRouteCrossConnectionsClient) GetCircuitMigrationInfoCallbackThenPoll(ctx context.Context, id ExpressRouteCrossConnectionId, input MigrateExpressRouteCircuitValidateAndHealthCheckRequest, callback func() error) error {
	result, err := c.GetCircuitMigrationInfo(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GetCircuitMigrationInfo: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetCircuitMigrationInfo: %+v", err)
	}

	return nil
}
