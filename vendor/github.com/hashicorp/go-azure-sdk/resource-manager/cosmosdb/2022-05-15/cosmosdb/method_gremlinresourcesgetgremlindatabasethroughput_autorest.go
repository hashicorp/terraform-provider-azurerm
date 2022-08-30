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

type GremlinResourcesGetGremlinDatabaseThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// GremlinResourcesGetGremlinDatabaseThroughput ...
func (c CosmosDBClient) GremlinResourcesGetGremlinDatabaseThroughput(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesGetGremlinDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForGremlinResourcesGetGremlinDatabaseThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGremlinResourcesGetGremlinDatabaseThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "GremlinResourcesGetGremlinDatabaseThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGremlinResourcesGetGremlinDatabaseThroughput prepares the GremlinResourcesGetGremlinDatabaseThroughput request.
func (c CosmosDBClient) preparerForGremlinResourcesGetGremlinDatabaseThroughput(ctx context.Context, id GremlinDatabaseId) (*http.Request, error) {
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

// responderForGremlinResourcesGetGremlinDatabaseThroughput handles the response to the GremlinResourcesGetGremlinDatabaseThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForGremlinResourcesGetGremlinDatabaseThroughput(resp *http.Response) (result GremlinResourcesGetGremlinDatabaseThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
