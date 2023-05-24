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

type DatabaseAccountsGetReadOnlyKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseAccountListReadOnlyKeysResult
}

// DatabaseAccountsGetReadOnlyKeys ...
func (c CosmosDBClient) DatabaseAccountsGetReadOnlyKeys(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsGetReadOnlyKeysOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsGetReadOnlyKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsGetReadOnlyKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsGetReadOnlyKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsGetReadOnlyKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsGetReadOnlyKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsGetReadOnlyKeys prepares the DatabaseAccountsGetReadOnlyKeys request.
func (c CosmosDBClient) preparerForDatabaseAccountsGetReadOnlyKeys(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/readonlykeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsGetReadOnlyKeys handles the response to the DatabaseAccountsGetReadOnlyKeys request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsGetReadOnlyKeys(resp *http.Response) (result DatabaseAccountsGetReadOnlyKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
