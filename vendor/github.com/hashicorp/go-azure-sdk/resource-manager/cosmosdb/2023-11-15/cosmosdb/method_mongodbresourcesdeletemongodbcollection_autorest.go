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

type MongoDBResourcesDeleteMongoDBCollectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesDeleteMongoDBCollection ...
func (c CosmosDBClient) MongoDBResourcesDeleteMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId) (result MongoDBResourcesDeleteMongoDBCollectionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesDeleteMongoDBCollection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesDeleteMongoDBCollection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesDeleteMongoDBCollection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesDeleteMongoDBCollection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesDeleteMongoDBCollectionThenPoll performs MongoDBResourcesDeleteMongoDBCollection then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesDeleteMongoDBCollectionThenPoll(ctx context.Context, id MongodbDatabaseCollectionId) error {
	result, err := c.MongoDBResourcesDeleteMongoDBCollection(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesDeleteMongoDBCollection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesDeleteMongoDBCollection: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesDeleteMongoDBCollection prepares the MongoDBResourcesDeleteMongoDBCollection request.
func (c CosmosDBClient) preparerForMongoDBResourcesDeleteMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId) (*http.Request, error) {
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

// senderForMongoDBResourcesDeleteMongoDBCollection sends the MongoDBResourcesDeleteMongoDBCollection request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesDeleteMongoDBCollection(ctx context.Context, req *http.Request) (future MongoDBResourcesDeleteMongoDBCollectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
