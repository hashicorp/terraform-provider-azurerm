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

type PercentileListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PercentileMetricListResult
}

type PercentileListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileListMetricsOperationOptions() PercentileListMetricsOperationOptions {
	return PercentileListMetricsOperationOptions{}
}

func (o PercentileListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PercentileListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// PercentileListMetrics ...
func (c CosmosDBClient) PercentileListMetrics(ctx context.Context, id DatabaseAccountId, options PercentileListMetricsOperationOptions) (result PercentileListMetricsOperationResponse, err error) {
	req, err := c.preparerForPercentileListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPercentileListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPercentileListMetrics prepares the PercentileListMetrics request.
func (c CosmosDBClient) preparerForPercentileListMetrics(ctx context.Context, id DatabaseAccountId, options PercentileListMetricsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/percentile/metrics", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPercentileListMetrics handles the response to the PercentileListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForPercentileListMetrics(resp *http.Response) (result PercentileListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
