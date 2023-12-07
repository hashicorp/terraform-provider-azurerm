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

type CollectionPartitionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PartitionMetricListResult
}

type CollectionPartitionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionListMetricsOperationOptions() CollectionPartitionListMetricsOperationOptions {
	return CollectionPartitionListMetricsOperationOptions{}
}

func (o CollectionPartitionListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionPartitionListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionPartitionListMetrics ...
func (c CosmosDBClient) CollectionPartitionListMetrics(ctx context.Context, id CollectionId, options CollectionPartitionListMetricsOperationOptions) (result CollectionPartitionListMetricsOperationResponse, err error) {
	req, err := c.preparerForCollectionPartitionListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionPartitionListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionPartitionListMetrics prepares the CollectionPartitionListMetrics request.
func (c CosmosDBClient) preparerForCollectionPartitionListMetrics(ctx context.Context, id CollectionId, options CollectionPartitionListMetricsOperationOptions) (*http.Request, error) {
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

// responderForCollectionPartitionListMetrics handles the response to the CollectionPartitionListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionPartitionListMetrics(resp *http.Response) (result CollectionPartitionListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
