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

type PercentileSourceTargetListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PercentileMetricListResult
}

type PercentileSourceTargetListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileSourceTargetListMetricsOperationOptions() PercentileSourceTargetListMetricsOperationOptions {
	return PercentileSourceTargetListMetricsOperationOptions{}
}

func (o PercentileSourceTargetListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PercentileSourceTargetListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// PercentileSourceTargetListMetrics ...
func (c CosmosDBClient) PercentileSourceTargetListMetrics(ctx context.Context, id SourceRegionTargetRegionId, options PercentileSourceTargetListMetricsOperationOptions) (result PercentileSourceTargetListMetricsOperationResponse, err error) {
	req, err := c.preparerForPercentileSourceTargetListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileSourceTargetListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileSourceTargetListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPercentileSourceTargetListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PercentileSourceTargetListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPercentileSourceTargetListMetrics prepares the PercentileSourceTargetListMetrics request.
func (c CosmosDBClient) preparerForPercentileSourceTargetListMetrics(ctx context.Context, id SourceRegionTargetRegionId, options PercentileSourceTargetListMetricsOperationOptions) (*http.Request, error) {
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

// responderForPercentileSourceTargetListMetrics handles the response to the PercentileSourceTargetListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForPercentileSourceTargetListMetrics(resp *http.Response) (result PercentileSourceTargetListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
