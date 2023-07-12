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

type DatabaseListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricListResult
}

type DatabaseListMetricsOperationOptions struct {
	Filter *string
}

func DefaultDatabaseListMetricsOperationOptions() DatabaseListMetricsOperationOptions {
	return DatabaseListMetricsOperationOptions{}
}

func (o DatabaseListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DatabaseListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// DatabaseListMetrics ...
func (c CosmosDBClient) DatabaseListMetrics(ctx context.Context, id DatabaseId, options DatabaseListMetricsOperationOptions) (result DatabaseListMetricsOperationResponse, err error) {
	req, err := c.preparerForDatabaseListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseListMetrics prepares the DatabaseListMetrics request.
func (c CosmosDBClient) preparerForDatabaseListMetrics(ctx context.Context, id DatabaseId, options DatabaseListMetricsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/metrics", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseListMetrics handles the response to the DatabaseListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseListMetrics(resp *http.Response) (result DatabaseListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
