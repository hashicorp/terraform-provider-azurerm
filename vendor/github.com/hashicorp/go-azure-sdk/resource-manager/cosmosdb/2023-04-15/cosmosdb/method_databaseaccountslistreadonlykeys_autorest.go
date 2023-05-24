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

type DatabaseAccountsListReadOnlyKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseAccountListReadOnlyKeysResult
}

// DatabaseAccountsListReadOnlyKeys ...
func (c CosmosDBClient) DatabaseAccountsListReadOnlyKeys(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsListReadOnlyKeysOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListReadOnlyKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListReadOnlyKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListReadOnlyKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListReadOnlyKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListReadOnlyKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListReadOnlyKeys prepares the DatabaseAccountsListReadOnlyKeys request.
func (c CosmosDBClient) preparerForDatabaseAccountsListReadOnlyKeys(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/readonlykeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsListReadOnlyKeys handles the response to the DatabaseAccountsListReadOnlyKeys request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListReadOnlyKeys(resp *http.Response) (result DatabaseAccountsListReadOnlyKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
