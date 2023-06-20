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

type GremlinResourcesMigrateGremlinDatabaseToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesMigrateGremlinDatabaseToAutoscale ...
func (c CosmosDBClient) GremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesMigrateGremlinDatabaseToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinDatabaseToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinDatabaseToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesMigrateGremlinDatabaseToAutoscaleThenPoll performs GremlinResourcesMigrateGremlinDatabaseToAutoscale then polls until it's completed
func (c CosmosDBClient) GremlinResourcesMigrateGremlinDatabaseToAutoscaleThenPoll(ctx context.Context, id GremlinDatabaseId) error {
	result, err := c.GremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesMigrateGremlinDatabaseToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesMigrateGremlinDatabaseToAutoscale: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesMigrateGremlinDatabaseToAutoscale prepares the GremlinResourcesMigrateGremlinDatabaseToAutoscale request.
func (c CosmosDBClient) preparerForGremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx context.Context, id GremlinDatabaseId) (*http.Request, error) {
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

// senderForGremlinResourcesMigrateGremlinDatabaseToAutoscale sends the GremlinResourcesMigrateGremlinDatabaseToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesMigrateGremlinDatabaseToAutoscale(ctx context.Context, req *http.Request) (future GremlinResourcesMigrateGremlinDatabaseToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
