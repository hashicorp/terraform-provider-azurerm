package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesMigrateGremlinGraphToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesMigrateGremlinGraphToAutoscale ...
func (c CosmosDBClient) GremlinResourcesMigrateGremlinGraphToAutoscale(ctx context.Context, id GraphId) (result GremlinResourcesMigrateGremlinGraphToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesMigrateGremlinGraphToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinGraphToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesMigrateGremlinGraphToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinGraphToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll performs GremlinResourcesMigrateGremlinGraphToAutoscale then polls until it's completed
func (c CosmosDBClient) GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll(ctx context.Context, id GraphId) error {
	result, err := c.GremlinResourcesMigrateGremlinGraphToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesMigrateGremlinGraphToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesMigrateGremlinGraphToAutoscale: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesMigrateGremlinGraphToAutoscale prepares the GremlinResourcesMigrateGremlinGraphToAutoscale request.
func (c CosmosDBClient) preparerForGremlinResourcesMigrateGremlinGraphToAutoscale(ctx context.Context, id GraphId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default/migrateToAutoscale", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGremlinResourcesMigrateGremlinGraphToAutoscale sends the GremlinResourcesMigrateGremlinGraphToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesMigrateGremlinGraphToAutoscale(ctx context.Context, req *http.Request) (future GremlinResourcesMigrateGremlinGraphToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
