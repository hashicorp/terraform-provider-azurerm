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

type CollectionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MetricListResult
}

type CollectionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionListMetricsOperationOptions() CollectionListMetricsOperationOptions {
	return CollectionListMetricsOperationOptions{}
}

func (o CollectionListMetricsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionListMetricsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionListMetrics ...
func (c CosmosDBClient) CollectionListMetrics(ctx context.Context, id CollectionId, options CollectionListMetricsOperationOptions) (result CollectionListMetricsOperationResponse, err error) {
	req, err := c.preparerForCollectionListMetrics(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetrics", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetrics", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionListMetrics(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListMetrics", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionListMetrics prepares the CollectionListMetrics request.
func (c CosmosDBClient) preparerForCollectionListMetrics(ctx context.Context, id CollectionId, options CollectionListMetricsOperationOptions) (*http.Request, error) {
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

// responderForCollectionListMetrics handles the response to the CollectionListMetrics request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionListMetrics(resp *http.Response) (result CollectionListMetricsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
