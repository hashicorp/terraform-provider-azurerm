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

type CollectionRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricListResult
}

type CollectionRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionRegionListMetricsOperationOptions() CollectionRegionListMetricsOperationOptions {
	return CollectionRegionListMetricsOperationOptions{}
}

func (o CollectionRegionListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionRegionListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionRegionListMetrics ...
func (c CosmosDBClient) CollectionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionRegionListMetricsOperationOptions) (result CollectionRegionListMetricsOperationResponse, err error) {
	req, err := c.preparerForCollectionRegionListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionRegionListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionRegionListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionRegionListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionRegionListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionRegionListMetrics prepares the CollectionRegionListMetrics request.
func (c CosmosDBClient) preparerForCollectionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionRegionListMetricsOperationOptions) (*http.Request, error) {
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

// responderForCollectionRegionListMetrics handles the response to the CollectionRegionListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionRegionListMetrics(resp *http.Response) (result CollectionRegionListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
