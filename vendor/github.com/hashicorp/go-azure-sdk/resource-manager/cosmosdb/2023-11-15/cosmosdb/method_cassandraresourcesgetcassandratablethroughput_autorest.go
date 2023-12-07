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

type CassandraResourcesGetCassandraTableThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// CassandraResourcesGetCassandraTableThroughput ...
func (c CosmosDBClient) CassandraResourcesGetCassandraTableThroughput(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesGetCassandraTableThroughputOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesGetCassandraTableThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTableThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTableThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesGetCassandraTableThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraTableThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesGetCassandraTableThroughput prepares the CassandraResourcesGetCassandraTableThroughput request.
func (c CosmosDBClient) preparerForCassandraResourcesGetCassandraTableThroughput(ctx context.Context, id CassandraKeyspaceTableId) (*http.Request, error) {
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

// responderForCassandraResourcesGetCassandraTableThroughput handles the response to the CassandraResourcesGetCassandraTableThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesGetCassandraTableThroughput(resp *http.Response) (result CassandraResourcesGetCassandraTableThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
