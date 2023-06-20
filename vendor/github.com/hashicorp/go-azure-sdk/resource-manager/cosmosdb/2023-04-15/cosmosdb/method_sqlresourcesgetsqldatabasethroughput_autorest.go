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

type SqlResourcesGetSqlDatabaseThroughputOperationResponse struct {
	HttpResponse *http.Response
	Model        *ThroughputSettingsGetResults
}

// SqlResourcesGetSqlDatabaseThroughput ...
func (c CosmosDBClient) SqlResourcesGetSqlDatabaseThroughput(ctx context.Context, id SqlDatabaseId) (result SqlResourcesGetSqlDatabaseThroughputOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesGetSqlDatabaseThroughput(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlDatabaseThroughput", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlDatabaseThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesGetSqlDatabaseThroughput(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlDatabaseThroughput", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesGetSqlDatabaseThroughput prepares the SqlResourcesGetSqlDatabaseThroughput request.
func (c CosmosDBClient) preparerForSqlResourcesGetSqlDatabaseThroughput(ctx context.Context, id SqlDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesGetSqlDatabaseThroughput handles the response to the SqlResourcesGetSqlDatabaseThroughput request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesGetSqlDatabaseThroughput(resp *http.Response) (result SqlResourcesGetSqlDatabaseThroughputOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
