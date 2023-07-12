package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesGetGremlinGraphOperationResponse struct {
	HttpResponse *http.Response
	Model        *GremlinGraphGetResults
}

// GremlinResourcesGetGremlinGraph ...
func (c CosmosDBClient) GremlinResourcesGetGremlinGraph(ctx context.Context, id GraphId) (result GremlinResourcesGetGremlinGraphOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesGetGremlinGraph(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraph", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraph", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesGetGremlinGraph(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraph", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesGetGremlinGraph prepares the GremlinResourcesGetGremlinGraph request.
func (c CosmosDBClient) preparerForGremlinResourcesGetGremlinGraph(ctx context.Context, id GraphId) (*http.Request, error) {
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

// responderForGremlinResourcesGetGremlinGraph handles the response to the GremlinResourcesGetGremlinGraph request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesGetGremlinGraph(resp *http.Response) (result GremlinResourcesGetGremlinGraphOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
