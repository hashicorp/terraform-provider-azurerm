package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesGetMongoDBDatabaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoDBDatabaseGetResults
}

// MongoDBResourcesGetMongoDBDatabase ...
func (c CosmosDBClient) MongoDBResourcesGetMongoDBDatabase(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesGetMongoDBDatabaseOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoDBDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabase", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoDBDatabase(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabase", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoDBDatabase prepares the MongoDBResourcesGetMongoDBDatabase request.
func (c CosmosDBClient) preparerForMongoDBResourcesGetMongoDBDatabase(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesGetMongoDBDatabase handles the response to the MongoDBResourcesGetMongoDBDatabase request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesGetMongoDBDatabase(resp *http.Response) (result MongoDBResourcesGetMongoDBDatabaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
