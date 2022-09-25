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

type GremlinResourcesGetGremlinGraphThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// GremlinResourcesGetGremlinGraphThroughput ...
func (c CosmosDBClient) GremlinResourcesGetGremlinGraphThroughput(ctx context.Context, id GraphId) (result GremlinResourcesGetGremlinGraphThroughputOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesGetGremlinGraphThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraphThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraphThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesGetGremlinGraphThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinGraphThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesGetGremlinGraphThroughput prepares the GremlinResourcesGetGremlinGraphThroughput request.
func (c CosmosDBClient) preparerForGremlinResourcesGetGremlinGraphThroughput(ctx context.Context, id GraphId) (*http.Request, error) {
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

// responderForGremlinResourcesGetGremlinGraphThroughput handles the response to the GremlinResourcesGetGremlinGraphThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesGetGremlinGraphThroughput(resp *http.Response) (result GremlinResourcesGetGremlinGraphThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
