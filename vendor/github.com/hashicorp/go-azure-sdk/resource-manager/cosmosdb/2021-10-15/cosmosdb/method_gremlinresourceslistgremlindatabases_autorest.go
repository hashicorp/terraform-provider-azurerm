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

type GremlinResourcesListGremlinDatabasesOperationResponse struct {
	HttpResponse *http.Response
	Model        *GremlinDatabaseListResult
}

// GremlinResourcesListGremlinDatabases ...
func (c CosmosDBClient) GremlinResourcesListGremlinDatabases(ctx context.Context, id DatabaseAccountId) (result GremlinResourcesListGremlinDatabasesOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesListGremlinDatabases(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinDatabases", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinDatabases", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesListGremlinDatabases(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinDatabases", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesListGremlinDatabases prepares the GremlinResourcesListGremlinDatabases request.
func (c CosmosDBClient) preparerForGremlinResourcesListGremlinDatabases(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/gremlinDatabases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGremlinResourcesListGremlinDatabases handles the response to the GremlinResourcesListGremlinDatabases request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesListGremlinDatabases(resp *http.Response) (result GremlinResourcesListGremlinDatabasesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
