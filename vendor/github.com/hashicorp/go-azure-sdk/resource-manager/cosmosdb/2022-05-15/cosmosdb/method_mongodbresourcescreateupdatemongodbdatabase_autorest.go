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

type MongoDBResourcesCreateUpdateMongoDBDatabaseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesCreateUpdateMongoDBDatabase ...
func (c CosmosDBClient) MongoDBResourcesCreateUpdateMongoDBDatabase(ctx context.Context, id MongodbDatabaseId, input MongoDBDatabaseCreateUpdateParameters) (result MongoDBResourcesCreateUpdateMongoDBDatabaseOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesCreateUpdateMongoDBDatabase(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesCreateUpdateMongoDBDatabase", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesCreateUpdateMongoDBDatabase(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesCreateUpdateMongoDBDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesCreateUpdateMongoDBDatabaseThenPoll performs MongoDBResourcesCreateUpdateMongoDBDatabase then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesCreateUpdateMongoDBDatabaseThenPoll(ctx context.Context, id MongodbDatabaseId, input MongoDBDatabaseCreateUpdateParameters) error {
	result, err := c.MongoDBResourcesCreateUpdateMongoDBDatabase(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesCreateUpdateMongoDBDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesCreateUpdateMongoDBDatabase: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesCreateUpdateMongoDBDatabase prepares the MongoDBResourcesCreateUpdateMongoDBDatabase request.
func (c CosmosDBClient) preparerForMongoDBResourcesCreateUpdateMongoDBDatabase(ctx context.Context, id MongodbDatabaseId, input MongoDBDatabaseCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMongoDBResourcesCreateUpdateMongoDBDatabase sends the MongoDBResourcesCreateUpdateMongoDBDatabase request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesCreateUpdateMongoDBDatabase(ctx context.Context, req *http.Request) (future MongoDBResourcesCreateUpdateMongoDBDatabaseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
