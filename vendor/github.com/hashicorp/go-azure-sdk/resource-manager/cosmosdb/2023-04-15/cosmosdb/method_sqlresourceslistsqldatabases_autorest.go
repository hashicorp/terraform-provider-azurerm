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

type SqlResourcesListSqlDatabasesOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlDatabaseListResult
}

// SqlResourcesListSqlDatabases ...
func (c CosmosDBClient) SqlResourcesListSqlDatabases(ctx context.Context, id DatabaseAccountId) (result SqlResourcesListSqlDatabasesOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesListSqlDatabases(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlDatabases", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlDatabases", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesListSqlDatabases(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlDatabases", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesListSqlDatabases prepares the SqlResourcesListSqlDatabases request.
func (c CosmosDBClient) preparerForSqlResourcesListSqlDatabases(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/sqlDatabases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesListSqlDatabases handles the response to the SqlResourcesListSqlDatabases request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesListSqlDatabases(resp *http.Response) (result SqlResourcesListSqlDatabasesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
