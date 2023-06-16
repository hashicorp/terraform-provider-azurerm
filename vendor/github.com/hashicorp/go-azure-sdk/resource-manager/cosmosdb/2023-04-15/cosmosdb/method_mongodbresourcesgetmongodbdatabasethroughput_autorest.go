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

type MongoDBResourcesGetMongoDBDatabaseThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// MongoDBResourcesGetMongoDBDatabaseThroughput ...
func (c CosmosDBClient) MongoDBResourcesGetMongoDBDatabaseThroughput(ctx context.Context, id MongodbDatabaseId) (result MongoDBResourcesGetMongoDBDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoDBDatabaseThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoDBDatabaseThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBDatabaseThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoDBDatabaseThroughput prepares the MongoDBResourcesGetMongoDBDatabaseThroughput request.
func (c CosmosDBClient) preparerForMongoDBResourcesGetMongoDBDatabaseThroughput(ctx context.Context, id MongodbDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesGetMongoDBDatabaseThroughput handles the response to the MongoDBResourcesGetMongoDBDatabaseThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesGetMongoDBDatabaseThroughput(resp *http.Response) (result MongoDBResourcesGetMongoDBDatabaseThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
