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

type DatabaseAccountsListUsagesOperationResponse struct {
	HttpResponse *http.Response
	Model        *UsagesResult
}

type DatabaseAccountsListUsagesOperationOptions struct {
	Filter *string
}

func DefaultDatabaseAccountsListUsagesOperationOptions() DatabaseAccountsListUsagesOperationOptions {
	return DatabaseAccountsListUsagesOperationOptions{}
}

func (o DatabaseAccountsListUsagesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DatabaseAccountsListUsagesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// DatabaseAccountsListUsages ...
func (c CosmosDBClient) DatabaseAccountsListUsages(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListUsagesOperationOptions) (result DatabaseAccountsListUsagesOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListUsages(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListUsages", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListUsages", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListUsages(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListUsages", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListUsages prepares the DatabaseAccountsListUsages request.
func (c CosmosDBClient) preparerForDatabaseAccountsListUsages(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListUsagesOperationOptions) (*http.Request, error) {
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

// responderForDatabaseAccountsListUsages handles the response to the DatabaseAccountsListUsages request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListUsages(resp *http.Response) (result DatabaseAccountsListUsagesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
