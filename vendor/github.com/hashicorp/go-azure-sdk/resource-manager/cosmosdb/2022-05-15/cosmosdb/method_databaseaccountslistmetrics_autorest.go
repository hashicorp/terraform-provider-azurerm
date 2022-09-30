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

type DatabaseAccountsListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricListResult
}

type DatabaseAccountsListMetricsOperationOptions struct {
	Filter *string
}

func DefaultDatabaseAccountsListMetricsOperationOptions() DatabaseAccountsListMetricsOperationOptions {
	return DatabaseAccountsListMetricsOperationOptions{}
}

func (o DatabaseAccountsListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DatabaseAccountsListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// DatabaseAccountsListMetrics ...
func (c CosmosDBClient) DatabaseAccountsListMetrics(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListMetricsOperationOptions) (result DatabaseAccountsListMetricsOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListMetrics prepares the DatabaseAccountsListMetrics request.
func (c CosmosDBClient) preparerForDatabaseAccountsListMetrics(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListMetricsOperationOptions) (*http.Request, error) {
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

// responderForDatabaseAccountsListMetrics handles the response to the DatabaseAccountsListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListMetrics(resp *http.Response) (result DatabaseAccountsListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
