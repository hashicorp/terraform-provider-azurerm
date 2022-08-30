package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesListMongoDBCollectionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoDBCollectionListResult
}

// MongoDBResourcesListMongoDBCollections ...
func (c CosmosDBClient) MongoDBResourcesListMongoDBCollections(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesListMongoDBCollectionsOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesListMongoDBCollections(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBCollections", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBCollections", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesListMongoDBCollections(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBCollections", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesListMongoDBCollections prepares the MongoDBResourcesListMongoDBCollections request.
func (c CosmosDBClient) preparerForMongoDBResourcesListMongoDBCollections(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/collections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesListMongoDBCollections handles the response to the MongoDBResourcesListMongoDBCollections request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesListMongoDBCollections(resp *http.Response) (result MongoDBResourcesListMongoDBCollectionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
