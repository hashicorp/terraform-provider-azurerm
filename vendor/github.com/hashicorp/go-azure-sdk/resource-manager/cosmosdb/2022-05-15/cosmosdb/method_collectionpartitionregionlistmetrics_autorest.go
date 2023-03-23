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

type CollectionPartitionRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PartitionMetricListResult
}

type CollectionPartitionRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionRegionListMetricsOperationOptions() CollectionPartitionRegionListMetricsOperationOptions {
	return CollectionPartitionRegionListMetricsOperationOptions{}
}

func (o CollectionPartitionRegionListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionPartitionRegionListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionPartitionRegionListMetrics ...
func (c CosmosDBClient) CollectionPartitionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions) (result CollectionPartitionRegionListMetricsOperationResponse, err error) {
	req, err := c.preparerForCollectionPartitionRegionListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionRegionListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionRegionListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionPartitionRegionListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionRegionListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionPartitionRegionListMetrics prepares the CollectionPartitionRegionListMetrics request.
func (c CosmosDBClient) preparerForCollectionPartitionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/partitions/metrics", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCollectionPartitionRegionListMetrics handles the response to the CollectionPartitionRegionListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionPartitionRegionListMetrics(resp *http.Response) (result CollectionPartitionRegionListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
