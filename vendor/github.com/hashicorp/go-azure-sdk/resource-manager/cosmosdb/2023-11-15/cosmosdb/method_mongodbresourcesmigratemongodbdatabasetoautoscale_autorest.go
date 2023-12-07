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

type MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesMigrateMongoDBDatabaseToAutoscale ...
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesMigrateMongoDBDatabaseToAutoscale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesMigrateMongoDBDatabaseToAutoscale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleThenPoll performs MongoDBResourcesMigrateMongoDBDatabaseToAutoscale then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleThenPoll(ctx context.Context, id MongodbDatabaseId) error {
	result, err := c.MongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesMigrateMongoDBDatabaseToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesMigrateMongoDBDatabaseToAutoscale: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale prepares the MongoDBResourcesMigrateMongoDBDatabaseToAutoscale request.
func (c CosmosDBClient) preparerForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
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

// senderForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale sends the MongoDBResourcesMigrateMongoDBDatabaseToAutoscale request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesMigrateMongoDBDatabaseToAutoscale(ctx context.Context, req *http.Request) (future MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
