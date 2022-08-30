package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraResourcesGetCassandraKeyspaceOperationResponse struct {
	HttpResponse *http.Response
	Model        *CassandraKeyspaceGetResults
}

// CassandraResourcesGetCassandraKeyspace ...
func (c CosmosDBClient) CassandraResourcesGetCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesGetCassandraKeyspaceOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesGetCassandraKeyspace(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspace", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspace", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesGetCassandraKeyspace(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspace", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesGetCassandraKeyspace prepares the CassandraResourcesGetCassandraKeyspace request.
func (c CosmosDBClient) preparerForCassandraResourcesGetCassandraKeyspace(ctx context.Context, id CassandraKeyspaceId) (*http.Request, error) {
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

// responderForCassandraResourcesGetCassandraKeyspace handles the response to the CassandraResourcesGetCassandraKeyspace request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesGetCassandraKeyspace(resp *http.Response) (result CassandraResourcesGetCassandraKeyspaceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
