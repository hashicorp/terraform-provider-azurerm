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

type SqlResourcesListSqlStoredProceduresOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlStoredProcedureListResult
}

// SqlResourcesListSqlStoredProcedures ...
func (c CosmosDBClient) SqlResourcesListSqlStoredProcedures(ctx context.Context, id ContainerId) (result SqlResourcesListSqlStoredProceduresOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesListSqlStoredProcedures(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlStoredProcedures", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlStoredProcedures", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesListSqlStoredProcedures(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlStoredProcedures", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesListSqlStoredProcedures prepares the SqlResourcesListSqlStoredProcedures request.
func (c CosmosDBClient) preparerForSqlResourcesListSqlStoredProcedures(ctx context.Context, id ContainerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/storedProcedures", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesListSqlStoredProcedures handles the response to the SqlResourcesListSqlStoredProcedures request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesListSqlStoredProcedures(resp *http.Response) (result SqlResourcesListSqlStoredProceduresOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
