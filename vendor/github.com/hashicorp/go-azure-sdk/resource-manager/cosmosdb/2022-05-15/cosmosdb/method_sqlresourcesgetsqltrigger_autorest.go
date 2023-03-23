package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesGetSqlTriggerOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlTriggerGetResults
}

// SqlResourcesGetSqlTrigger ...
func (c CosmosDBClient) SqlResourcesGetSqlTrigger(ctx context.Context, id TriggerId) (result SqlResourcesGetSqlTriggerOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesGetSqlTrigger(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlTrigger", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlTrigger", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesGetSqlTrigger(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlTrigger", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesGetSqlTrigger prepares the SqlResourcesGetSqlTrigger request.
func (c CosmosDBClient) preparerForSqlResourcesGetSqlTrigger(ctx context.Context, id TriggerId) (*http.Request, error) {
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

// responderForSqlResourcesGetSqlTrigger handles the response to the SqlResourcesGetSqlTrigger request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesGetSqlTrigger(resp *http.Response) (result SqlResourcesGetSqlTriggerOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
