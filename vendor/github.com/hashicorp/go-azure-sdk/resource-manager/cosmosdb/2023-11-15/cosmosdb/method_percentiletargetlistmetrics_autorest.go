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

type PercentileTargetListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PercentileMetricListResult
}

type PercentileTargetListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileTargetListMetricsOperationOptions() PercentileTargetListMetricsOperationOptions {
	return PercentileTargetListMetricsOperationOptions{}
}

func (o PercentileTargetListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PercentileTargetListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// PercentileTargetListMetrics ...
func (c CosmosDBClient) PercentileTargetListMetrics(ctx context.Context, id TargetRegionId, options PercentileTargetListMetricsOperationOptions) (result PercentileTargetListMetricsOperationResponse, err error) {
	req, err := c.preparerForPercentileTargetListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileTargetListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileTargetListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPercentileTargetListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileTargetListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPercentileTargetListMetrics prepares the PercentileTargetListMetrics request.
func (c CosmosDBClient) preparerForPercentileTargetListMetrics(ctx context.Context, id TargetRegionId, options PercentileTargetListMetricsOperationOptions) (*http.Request, error) {
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

// responderForPercentileTargetListMetrics handles the response to the PercentileTargetListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForPercentileTargetListMetrics(resp *http.Response) (result PercentileTargetListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
