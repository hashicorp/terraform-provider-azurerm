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

type DatabaseListUsagesOperationResponse struct {
	HttpResponse *http.Response
	Model        *UsagesResult
}

type DatabaseListUsagesOperationOptions struct {
	Filter *string
}

func DefaultDatabaseListUsagesOperationOptions() DatabaseListUsagesOperationOptions {
	return DatabaseListUsagesOperationOptions{}
}

func (o DatabaseListUsagesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DatabaseListUsagesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// DatabaseListUsages ...
func (c CosmosDBClient) DatabaseListUsages(ctx context.Context, id DatabaseId, options DatabaseListUsagesOperationOptions) (result DatabaseListUsagesOperationResponse, err error) {
	req, err := c.preparerForDatabaseListUsages(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListUsages", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListUsages", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseListUsages(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseListUsages", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseListUsages prepares the DatabaseListUsages request.
func (c CosmosDBClient) preparerForDatabaseListUsages(ctx context.Context, id DatabaseId, options DatabaseListUsagesOperationOptions) (*http.Request, error) {
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

// responderForDatabaseListUsages handles the response to the DatabaseListUsages request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseListUsages(resp *http.Response) (result DatabaseListUsagesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
