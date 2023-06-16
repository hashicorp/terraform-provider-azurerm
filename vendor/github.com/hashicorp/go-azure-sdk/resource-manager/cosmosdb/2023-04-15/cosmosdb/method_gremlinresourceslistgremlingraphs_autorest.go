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

type GremlinResourcesListGremlinGraphsOperationResponse struct {
	HttpResponse *http.Response
	Model        *GremlinGraphListResult
}

// GremlinResourcesListGremlinGraphs ...
func (c CosmosDBClient) GremlinResourcesListGremlinGraphs(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesListGremlinGraphsOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesListGremlinGraphs(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinGraphs", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinGraphs", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesListGremlinGraphs(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesListGremlinGraphs", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesListGremlinGraphs prepares the GremlinResourcesListGremlinGraphs request.
func (c CosmosDBClient) preparerForGremlinResourcesListGremlinGraphs(ctx context.Context, id GremlinDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/graphs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGremlinResourcesListGremlinGraphs handles the response to the GremlinResourcesListGremlinGraphs request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesListGremlinGraphs(resp *http.Response) (result GremlinResourcesListGremlinGraphsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
