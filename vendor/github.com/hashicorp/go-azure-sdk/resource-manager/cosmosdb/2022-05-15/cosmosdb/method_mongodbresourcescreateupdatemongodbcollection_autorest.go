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

type MongoDBResourcesCreateUpdateMongoDBCollectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesCreateUpdateMongoDBCollection ...
func (c CosmosDBClient) MongoDBResourcesCreateUpdateMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId, input MongoDBCollectionCreateUpdateParameters) (result MongoDBResourcesCreateUpdateMongoDBCollectionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesCreateUpdateMongoDBCollection(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesCreateUpdateMongoDBCollection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesCreateUpdateMongoDBCollection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesCreateUpdateMongoDBCollection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll performs MongoDBResourcesCreateUpdateMongoDBCollection then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll(ctx context.Context, id MongodbDatabaseCollectionId, input MongoDBCollectionCreateUpdateParameters) error {
	result, err := c.MongoDBResourcesCreateUpdateMongoDBCollection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesCreateUpdateMongoDBCollection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesCreateUpdateMongoDBCollection: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesCreateUpdateMongoDBCollection prepares the MongoDBResourcesCreateUpdateMongoDBCollection request.
func (c CosmosDBClient) preparerForMongoDBResourcesCreateUpdateMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId, input MongoDBCollectionCreateUpdateParameters) (*http.Request, error) {
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

// senderForMongoDBResourcesCreateUpdateMongoDBCollection sends the MongoDBResourcesCreateUpdateMongoDBCollection request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForMongoDBResourcesCreateUpdateMongoDBCollection(ctx context.Context, req *http.Request) (future MongoDBResourcesCreateUpdateMongoDBCollectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
