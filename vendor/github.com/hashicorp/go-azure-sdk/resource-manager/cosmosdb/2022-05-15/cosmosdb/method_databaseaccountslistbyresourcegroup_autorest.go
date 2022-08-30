package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseAccountsListResult
}

// DatabaseAccountsListByResourceGroup ...
func (c CosmosDBClient) DatabaseAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result DatabaseAccountsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseAccountsListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseAccountsListByResourceGroup prepares the DatabaseAccountsListByResourceGroup request.
func (c CosmosDBClient) preparerForDatabaseAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DocumentDB/databaseAccounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseAccountsListByResourceGroup handles the response to the DatabaseAccountsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForDatabaseAccountsListByResourceGroup(resp *http.Response) (result DatabaseAccountsListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
