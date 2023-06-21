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

type MongoDBResourcesDeleteMongoDBDatabaseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesDeleteMongoDBDatabase ...
func (c CosmosDBClient) MongoDBResourcesDeleteMongoDBDatabase(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesDeleteMongoDBDatabaseOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesDeleteMongoDBDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesDeleteMongoDBDatabase", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesDeleteMongoDBDatabase(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesDeleteMongoDBDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesDeleteMongoDBDatabaseThenPoll performs MongoDBResourcesDeleteMongoDBDatabase then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesDeleteMongoDBDatabaseThenPoll(ctx context.Context, id MongodbDatabaseId) error {
	result, err := c.MongoDBResourcesDeleteMongoDBDatabase(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesDeleteMongoDBDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesDeleteMongoDBDatabase: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesDeleteMongoDBDatabase prepares the MongoDBResourcesDeleteMongoDBDatabase request.
func (c CosmosDBClient) preparerForMongoDBResourcesDeleteMongoDBDatabase(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMongoDBResourcesDeleteMongoDBDatabase sends the MongoDBResourcesDeleteMongoDBDatabase request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesDeleteMongoDBDatabase(ctx context.Context, req *http.Request) (future MongoDBResourcesDeleteMongoDBDatabaseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
