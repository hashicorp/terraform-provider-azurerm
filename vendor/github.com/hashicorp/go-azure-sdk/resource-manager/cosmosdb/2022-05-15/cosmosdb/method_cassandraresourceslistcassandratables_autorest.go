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

type CassandraResourcesListCassandraTablesOperationResponse struct {
	HttpResponse *http.Response
	Model        *CassandraTableListResult
}

// CassandraResourcesListCassandraTables ...
func (c CosmosDBClient) CassandraResourcesListCassandraTables(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesListCassandraTablesOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesListCassandraTables(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraTables", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraTables", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesListCassandraTables(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraTables", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesListCassandraTables prepares the CassandraResourcesListCassandraTables request.
func (c CosmosDBClient) preparerForCassandraResourcesListCassandraTables(ctx context.Context, id CassandraKeyspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tables", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCassandraResourcesListCassandraTables handles the response to the CassandraResourcesListCassandraTables request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesListCassandraTables(resp *http.Response) (result CassandraResourcesListCassandraTablesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
