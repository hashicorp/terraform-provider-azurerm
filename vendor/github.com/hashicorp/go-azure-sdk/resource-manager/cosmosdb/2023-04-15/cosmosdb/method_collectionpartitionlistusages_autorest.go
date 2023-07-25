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

type CollectionPartitionListUsagesOperationResponse struct {
	HttpResponse *http.Response
	Model        *PartitionUsagesResult
}

type CollectionPartitionListUsagesOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionListUsagesOperationOptions() CollectionPartitionListUsagesOperationOptions {
	return CollectionPartitionListUsagesOperationOptions{}
}

func (o CollectionPartitionListUsagesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionPartitionListUsagesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionPartitionListUsages ...
func (c CosmosDBClient) CollectionPartitionListUsages(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions) (result CollectionPartitionListUsagesOperationResponse, err error) {
	req, err := c.preparerForCollectionPartitionListUsages(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListUsages", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListUsages", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionPartitionListUsages(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionPartitionListUsages", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionPartitionListUsages prepares the CollectionPartitionListUsages request.
func (c CosmosDBClient) preparerForCollectionPartitionListUsages(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/partitions/usages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCollectionPartitionListUsages handles the response to the CollectionPartitionListUsages request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionPartitionListUsages(resp *http.Response) (result CollectionPartitionListUsagesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
