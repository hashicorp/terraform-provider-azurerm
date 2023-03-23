package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesGetSqlUserDefinedFunctionOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlUserDefinedFunctionGetResults
}

// SqlResourcesGetSqlUserDefinedFunction ...
func (c CosmosDBClient) SqlResourcesGetSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId) (result SqlResourcesGetSqlUserDefinedFunctionOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesGetSqlUserDefinedFunction(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlUserDefinedFunction", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlUserDefinedFunction", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesGetSqlUserDefinedFunction(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesGetSqlUserDefinedFunction", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesGetSqlUserDefinedFunction prepares the SqlResourcesGetSqlUserDefinedFunction request.
func (c CosmosDBClient) preparerForSqlResourcesGetSqlUserDefinedFunction(ctx context.Context, id UserDefinedFunctionId) (*http.Request, error) {
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

// responderForSqlResourcesGetSqlUserDefinedFunction handles the response to the SqlResourcesGetSqlUserDefinedFunction request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesGetSqlUserDefinedFunction(resp *http.Response) (result SqlResourcesGetSqlUserDefinedFunctionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
