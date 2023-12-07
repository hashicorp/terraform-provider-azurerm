package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraResourcesGetCassandraTableOperationResponse struct {
	HttpResponse *http.Response
	Model        *CassandraTableGetResults
}

// CassandraResourcesGetCassandraTable ...
func (c CosmosDBClient) CassandraResourcesGetCassandraTable(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesGetCassandraTableOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesGetCassandraTable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesGetCassandraTable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesGetCassandraTable prepares the CassandraResourcesGetCassandraTable request.
func (c CosmosDBClient) preparerForCassandraResourcesGetCassandraTable(ctx context.Context, id CassandraKeyspaceTableId) (*http.Request, error) {
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

// responderForCassandraResourcesGetCassandraTable handles the response to the CassandraResourcesGetCassandraTable request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesGetCassandraTable(resp *http.Response) (result CassandraResourcesGetCassandraTableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
