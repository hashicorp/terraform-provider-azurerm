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

type GremlinResourcesMigrateGremlinDatabaseToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GremlinResourcesMigrateGremlinDatabaseToManualThroughput ...
func (c CosmosDBClient) GremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesMigrateGremlinDatabaseToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinDatabaseToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesMigrateGremlinDatabaseToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GremlinResourcesMigrateGremlinDatabaseToManualThroughputThenPoll performs GremlinResourcesMigrateGremlinDatabaseToManualThroughput then polls until it's completed
func (c CosmosDBClient) GremlinResourcesMigrateGremlinDatabaseToManualThroughputThenPoll(ctx context.Context, id GremlinDatabaseId) error {
	result, err := c.GremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesMigrateGremlinDatabaseToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GremlinResourcesMigrateGremlinDatabaseToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForGremlinResourcesMigrateGremlinDatabaseToManualThroughput prepares the GremlinResourcesMigrateGremlinDatabaseToManualThroughput request.
func (c CosmosDBClient) preparerForGremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx context.Context, id GremlinDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default/migrateToManualThroughput", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGremlinResourcesMigrateGremlinDatabaseToManualThroughput sends the GremlinResourcesMigrateGremlinDatabaseToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForGremlinResourcesMigrateGremlinDatabaseToManualThroughput(ctx context.Context, req *http.Request) (future GremlinResourcesMigrateGremlinDatabaseToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
