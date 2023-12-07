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

type PartitionKeyRangeIdRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PartitionMetricListResult
}

type PartitionKeyRangeIdRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPartitionKeyRangeIdRegionListMetricsOperationOptions() PartitionKeyRangeIdRegionListMetricsOperationOptions {
	return PartitionKeyRangeIdRegionListMetricsOperationOptions{}
}

func (o PartitionKeyRangeIdRegionListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PartitionKeyRangeIdRegionListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// PartitionKeyRangeIdRegionListMetrics ...
func (c CosmosDBClient) PartitionKeyRangeIdRegionListMetrics(ctx context.Context, id CollectionPartitionKeyRangeIdId, options PartitionKeyRangeIdRegionListMetricsOperationOptions) (result PartitionKeyRangeIdRegionListMetricsOperationResponse, err error) {
	req, err := c.preparerForPartitionKeyRangeIdRegionListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdRegionListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdRegionListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPartitionKeyRangeIdRegionListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "PartitionKeyRangeIdRegionListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPartitionKeyRangeIdRegionListMetrics prepares the PartitionKeyRangeIdRegionListMetrics request.
func (c CosmosDBClient) preparerForPartitionKeyRangeIdRegionListMetrics(ctx context.Context, id CollectionPartitionKeyRangeIdId, options PartitionKeyRangeIdRegionListMetricsOperationOptions) (*http.Request, error) {
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

// responderForPartitionKeyRangeIdRegionListMetrics handles the response to the PartitionKeyRangeIdRegionListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForPartitionKeyRangeIdRegionListMetrics(resp *http.Response) (result PartitionKeyRangeIdRegionListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
