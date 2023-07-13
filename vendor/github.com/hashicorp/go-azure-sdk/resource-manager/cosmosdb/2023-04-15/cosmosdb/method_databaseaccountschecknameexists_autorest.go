package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountsCheckNameExistsOperationResponse struct {
	HttpResponse *http.Response
}

// DatabaseAccountsCheckNameExists ...
func (c CosmosDBClient) DatabaseAccountsCheckNameExists(ctx context.Context, id DatabaseAccountNameId) (result DatabaseAccountsCheckNameExistsOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsCheckNameExists(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsCheckNameExists", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsCheckNameExists", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsCheckNameExists(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsCheckNameExists", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsCheckNameExists prepares the DatabaseAccountsCheckNameExists request.
func (c CosmosDBClient) preparerForDatabaseAccountsCheckNameExists(ctx context.Context, id DatabaseAccountNameId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsHead(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsCheckNameExists handles the response to the DatabaseAccountsCheckNameExists request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsCheckNameExists(resp *http.Response) (result DatabaseAccountsCheckNameExistsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
