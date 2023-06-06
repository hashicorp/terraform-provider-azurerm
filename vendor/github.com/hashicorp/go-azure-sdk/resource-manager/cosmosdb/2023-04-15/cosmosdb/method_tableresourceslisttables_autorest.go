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

type TableResourcesListTablesOperationResponse struct {
	HttpResponse *http.Response
	Model        *TableListResult
}

// TableResourcesListTables ...
func (c CosmosDBClient) TableResourcesListTables(ctx context.Context, id DatabaseAccountId) (result TableResourcesListTablesOperationResponse, err error) {
	req, err := c.preparerForTableResourcesListTables(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesListTables", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesListTables", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableResourcesListTables(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesListTables", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableResourcesListTables prepares the TableResourcesListTables request.
func (c CosmosDBClient) preparerForTableResourcesListTables(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tables", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableResourcesListTables handles the response to the TableResourcesListTables request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForTableResourcesListTables(resp *http.Response) (result TableResourcesListTablesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
