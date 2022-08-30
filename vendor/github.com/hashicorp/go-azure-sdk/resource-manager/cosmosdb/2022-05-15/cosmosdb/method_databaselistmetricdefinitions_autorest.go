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

type DatabaseListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricDefinitionsListResult
}

// DatabaseListMetricDefinitions ...
func (c CosmosDBClient) DatabaseListMetricDefinitions(ctx context.Context, id DatabaseId) (result DatabaseListMetricDefinitionsOperationResponse, err error) {
	req, err := c.preparerForDatabaseListMetricDefinitions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetricDefinitions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetricDefinitions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseListMetricDefinitions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetricDefinitions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseListMetricDefinitions prepares the DatabaseListMetricDefinitions request.
func (c CosmosDBClient) preparerForDatabaseListMetricDefinitions(ctx context.Context, id DatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/metricDefinitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseListMetricDefinitions handles the response to the DatabaseListMetricDefinitions request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseListMetricDefinitions(resp *http.Response) (result DatabaseListMetricDefinitionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
