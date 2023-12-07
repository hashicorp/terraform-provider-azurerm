package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesGetGremlinDatabaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *GremlinDatabaseGetResults
}

// GremlinResourcesGetGremlinDatabase ...
func (c CosmosDBClient) GremlinResourcesGetGremlinDatabase(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesGetGremlinDatabaseOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesGetGremlinDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabase", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesGetGremlinDatabase(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabase", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesGetGremlinDatabase prepares the GremlinResourcesGetGremlinDatabase request.
func (c CosmosDBClient) preparerForGremlinResourcesGetGremlinDatabase(ctx context.Context, id GremlinDatabaseId) (*http.Request, error) {
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

// responderForGremlinResourcesGetGremlinDatabase handles the response to the GremlinResourcesGetGremlinDatabase request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesGetGremlinDatabase(resp *http.Response) (result GremlinResourcesGetGremlinDatabaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
