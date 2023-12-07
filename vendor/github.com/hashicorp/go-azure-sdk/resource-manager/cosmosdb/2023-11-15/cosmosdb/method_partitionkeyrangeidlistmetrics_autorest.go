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

type PartitionKeyRangeIdListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PartitionMetricListResult
}

type PartitionKeyRangeIdListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPartitionKeyRangeIdListMetricsOperationOptions() PartitionKeyRangeIdListMetricsOperationOptions {
	return PartitionKeyRangeIdListMetricsOperationOptions{}
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// PartitionKeyRangeIdListMetrics ...
func (c CosmosDBClient) PartitionKeyRangeIdListMetrics(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions) (result PartitionKeyRangeIdListMetricsOperationResponse, err error) {
	req, err := c.preparerForPartitionKeyRangeIdListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPartitionKeyRangeIdListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPartitionKeyRangeIdListMetrics prepares the PartitionKeyRangeIdListMetrics request.
func (c CosmosDBClient) preparerForPartitionKeyRangeIdListMetrics(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions) (*http.Request, error) {
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

// responderForPartitionKeyRangeIdListMetrics handles the response to the PartitionKeyRangeIdListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForPartitionKeyRangeIdListMetrics(resp *http.Response) (result PartitionKeyRangeIdListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
