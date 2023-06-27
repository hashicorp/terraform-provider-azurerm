package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesGetSqlStoredProcedureOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlStoredProcedureGetResults
}

// SqlResourcesGetSqlStoredProcedure ...
func (c CosmosDBClient) SqlResourcesGetSqlStoredProcedure(ctx context.Context, id StoredProcedureId) (result SqlResourcesGetSqlStoredProcedureOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesGetSqlStoredProcedure(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlStoredProcedure", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlStoredProcedure", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesGetSqlStoredProcedure(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlStoredProcedure", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesGetSqlStoredProcedure prepares the SqlResourcesGetSqlStoredProcedure request.
func (c CosmosDBClient) preparerForSqlResourcesGetSqlStoredProcedure(ctx context.Context, id StoredProcedureId) (*http.Request, error) {
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

// responderForSqlResourcesGetSqlStoredProcedure handles the response to the SqlResourcesGetSqlStoredProcedure request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesGetSqlStoredProcedure(resp *http.Response) (result SqlResourcesGetSqlStoredProcedureOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
