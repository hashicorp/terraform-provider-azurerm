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

type DatabaseAccountsListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseAccountListKeysResult
}

// DatabaseAccountsListKeys ...
func (c CosmosDBClient) DatabaseAccountsListKeys(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsListKeysOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListKeys prepares the DatabaseAccountsListKeys request.
func (c CosmosDBClient) preparerForDatabaseAccountsListKeys(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsListKeys handles the response to the DatabaseAccountsListKeys request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListKeys(resp *http.Response) (result DatabaseAccountsListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
