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

type MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput ...
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputThenPoll performs MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputThenPoll(ctx context.Context, id MongodbDatabaseId) error {
	result, err := c.MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput prepares the MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput request.
func (c CosmosDBClient) preparerForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
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

// senderForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput sends the MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesMigrateMongoDBDatabaseToManualThroughput(ctx context.Context, req *http.Request) (future MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
