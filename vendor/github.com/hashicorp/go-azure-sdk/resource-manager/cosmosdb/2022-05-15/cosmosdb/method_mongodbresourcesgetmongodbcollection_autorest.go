package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesGetMongoDBCollectionOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoDBCollectionGetResults
}

// MongoDBResourcesGetMongoDBCollection ...
func (c CosmosDBClient) MongoDBResourcesGetMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId) (result MongoDBResourcesGetMongoDBCollectionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoDBCollection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollection", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollection", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoDBCollection(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "MongoDBResourcesGetMongoDBCollection", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoDBCollection prepares the MongoDBResourcesGetMongoDBCollection request.
func (c CosmosDBClient) preparerForMongoDBResourcesGetMongoDBCollection(ctx context.Context, id MongodbDatabaseCollectionId) (*http.Request, error) {
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

// responderForMongoDBResourcesGetMongoDBCollection handles the response to the MongoDBResourcesGetMongoDBCollection request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForMongoDBResourcesGetMongoDBCollection(resp *http.Response) (result MongoDBResourcesGetMongoDBCollectionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
