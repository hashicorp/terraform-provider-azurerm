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

type CollectionListUsagesOperationResponse struct {
	HttpResponse *http.Response
	Model        *UsagesResult
}

type CollectionListUsagesOperationOptions struct {
	Filter *string
}

func DefaultCollectionListUsagesOperationOptions() CollectionListUsagesOperationOptions {
	return CollectionListUsagesOperationOptions{}
}

func (o CollectionListUsagesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CollectionListUsagesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// CollectionListUsages ...
func (c CosmosDBClient) CollectionListUsages(ctx context.Context, id CollectionId, options CollectionListUsagesOperationOptions) (result CollectionListUsagesOperationResponse, err error) {
	req, err := c.preparerForCollectionListUsages(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListUsages", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListUsages", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCollectionListUsages(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "CollectionListUsages", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCollectionListUsages prepares the CollectionListUsages request.
func (c CosmosDBClient) preparerForCollectionListUsages(ctx context.Context, id CollectionId, options CollectionListUsagesOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/usages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCollectionListUsages handles the response to the CollectionListUsages request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForCollectionListUsages(resp *http.Response) (result CollectionListUsagesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
