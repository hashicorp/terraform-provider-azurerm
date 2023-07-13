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

type DatabaseAccountsListConnectionStringsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseAccountListConnectionStringsResult
}

// DatabaseAccountsListConnectionStrings ...
func (c CosmosDBClient) DatabaseAccountsListConnectionStrings(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsListConnectionStringsOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListConnectionStrings(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListConnectionStrings", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListConnectionStrings", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListConnectionStrings(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListConnectionStrings", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListConnectionStrings prepares the DatabaseAccountsListConnectionStrings request.
func (c CosmosDBClient) preparerForDatabaseAccountsListConnectionStrings(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listConnectionStrings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsListConnectionStrings handles the response to the DatabaseAccountsListConnectionStrings request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListConnectionStrings(resp *http.Response) (result DatabaseAccountsListConnectionStringsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
