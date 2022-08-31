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

type MongoDBResourcesGetMongoDBCollectionThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// MongoDBResourcesGetMongoDBCollectionThroughput ...
func (c CosmosDBClient) MongoDBResourcesGetMongoDBCollectionThroughput(ctx context.Context, id MongodbDatabaseCollectionId) (result MongoDBResourcesGetMongoDBCollectionThroughputOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoDBCollectionThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollectionThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollectionThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoDBCollectionThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollectionThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoDBCollectionThroughput prepares the MongoDBResourcesGetMongoDBCollectionThroughput request.
func (c CosmosDBClient) preparerForMongoDBResourcesGetMongoDBCollectionThroughput(ctx context.Context, id MongodbDatabaseCollectionId) (*http.Request, error) {
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

// responderForMongoDBResourcesGetMongoDBCollectionThroughput handles the response to the MongoDBResourcesGetMongoDBCollectionThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesGetMongoDBCollectionThroughput(resp *http.Response) (result MongoDBResourcesGetMongoDBCollectionThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
