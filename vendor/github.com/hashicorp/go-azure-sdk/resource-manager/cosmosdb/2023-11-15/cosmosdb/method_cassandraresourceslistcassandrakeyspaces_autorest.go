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

type CassandraResourcesListCassandraKeyspacesOperationResponse struct {
	HttpResponse *http.Response
	Model        *CassandraKeyspaceListResult
}

// CassandraResourcesListCassandraKeyspaces ...
func (c CosmosDBClient) CassandraResourcesListCassandraKeyspaces(ctx context.Context, id DatabaseAccountId) (result CassandraResourcesListCassandraKeyspacesOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesListCassandraKeyspaces(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraKeyspaces", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraKeyspaces", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesListCassandraKeyspaces(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesListCassandraKeyspaces", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesListCassandraKeyspaces prepares the CassandraResourcesListCassandraKeyspaces request.
func (c CosmosDBClient) preparerForCassandraResourcesListCassandraKeyspaces(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cassandraKeyspaces", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCassandraResourcesListCassandraKeyspaces handles the response to the CassandraResourcesListCassandraKeyspaces request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesListCassandraKeyspaces(resp *http.Response) (result CassandraResourcesListCassandraKeyspacesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
