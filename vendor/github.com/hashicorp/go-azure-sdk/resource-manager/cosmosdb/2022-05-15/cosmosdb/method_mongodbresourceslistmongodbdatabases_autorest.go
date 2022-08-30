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

type MongoDBResourcesListMongoDBDatabasesOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoDBDatabaseListResult
}

// MongoDBResourcesListMongoDBDatabases ...
func (c CosmosDBClient) MongoDBResourcesListMongoDBDatabases(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoDBDatabasesOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesListMongoDBDatabases(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBDatabases", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBDatabases", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesListMongoDBDatabases(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesListMongoDBDatabases", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesListMongoDBDatabases prepares the MongoDBResourcesListMongoDBDatabases request.
func (c CosmosDBClient) preparerForMongoDBResourcesListMongoDBDatabases(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/mongodbDatabases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesListMongoDBDatabases handles the response to the MongoDBResourcesListMongoDBDatabases request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesListMongoDBDatabases(resp *http.Response) (result MongoDBResourcesListMongoDBDatabasesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
