package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableResourcesGetTableOperationResponse struct {
	HttpResponse *http.Response
	Model        *TableGetResults
}

// TableResourcesGetTable ...
func (c CosmosDBClient) TableResourcesGetTable(ctx context.Context, id TableId) (result TableResourcesGetTableOperationResponse, err error) {
	req, err := c.preparerForTableResourcesGetTable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableResourcesGetTable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesGetTable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableResourcesGetTable prepares the TableResourcesGetTable request.
func (c CosmosDBClient) preparerForTableResourcesGetTable(ctx context.Context, id TableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableResourcesGetTable handles the response to the TableResourcesGetTable request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForTableResourcesGetTable(resp *http.Response) (result TableResourcesGetTableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
