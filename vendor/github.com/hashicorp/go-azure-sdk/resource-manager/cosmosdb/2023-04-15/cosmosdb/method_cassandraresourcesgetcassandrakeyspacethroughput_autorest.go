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

type CassandraResourcesGetCassandraKeyspaceThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// CassandraResourcesGetCassandraKeyspaceThroughput ...
func (c CosmosDBClient) CassandraResourcesGetCassandraKeyspaceThroughput(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesGetCassandraKeyspaceThroughputOperationResponse, err error) {
	req, err := c.preparerForCassandraResourcesGetCassandraKeyspaceThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspaceThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspaceThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCassandraResourcesGetCassandraKeyspaceThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CassandraResourcesGetCassandraKeyspaceThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCassandraResourcesGetCassandraKeyspaceThroughput prepares the CassandraResourcesGetCassandraKeyspaceThroughput request.
func (c CosmosDBClient) preparerForCassandraResourcesGetCassandraKeyspaceThroughput(ctx context.Context, id CassandraKeyspaceId) (*http.Request, error) {
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

// responderForCassandraResourcesGetCassandraKeyspaceThroughput handles the response to the CassandraResourcesGetCassandraKeyspaceThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCassandraResourcesGetCassandraKeyspaceThroughput(resp *http.Response) (result CassandraResourcesGetCassandraKeyspaceThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
