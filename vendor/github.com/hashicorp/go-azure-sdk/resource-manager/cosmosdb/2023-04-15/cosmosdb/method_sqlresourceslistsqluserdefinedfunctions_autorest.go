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

type SqlResourcesListSqlUserDefinedFunctionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlUserDefinedFunctionListResult
}

// SqlResourcesListSqlUserDefinedFunctions ...
func (c CosmosDBClient) SqlResourcesListSqlUserDefinedFunctions(ctx context.Context, id ContainerId) (result SqlResourcesListSqlUserDefinedFunctionsOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesListSqlUserDefinedFunctions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlUserDefinedFunctions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlUserDefinedFunctions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesListSqlUserDefinedFunctions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlUserDefinedFunctions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesListSqlUserDefinedFunctions prepares the SqlResourcesListSqlUserDefinedFunctions request.
func (c CosmosDBClient) preparerForSqlResourcesListSqlUserDefinedFunctions(ctx context.Context, id ContainerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/userDefinedFunctions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesListSqlUserDefinedFunctions handles the response to the SqlResourcesListSqlUserDefinedFunctions request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesListSqlUserDefinedFunctions(resp *http.Response) (result SqlResourcesListSqlUserDefinedFunctionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
