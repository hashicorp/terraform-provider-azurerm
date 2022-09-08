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

type SqlResourcesListSqlContainersOperationResponse struct {
	HttpResponse *http.Response
	Model        *SqlContainerListResult
}

// SqlResourcesListSqlContainers ...
func (c CosmosDBClient) SqlResourcesListSqlContainers(ctx context.Context, id SqlDatabaseId) (result SqlResourcesListSqlContainersOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesListSqlContainers(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlContainers", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlContainers", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSqlResourcesListSqlContainers(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesListSqlContainers", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSqlResourcesListSqlContainers prepares the SqlResourcesListSqlContainers request.
func (c CosmosDBClient) preparerForSqlResourcesListSqlContainers(ctx context.Context, id SqlDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/containers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSqlResourcesListSqlContainers handles the response to the SqlResourcesListSqlContainers request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForSqlResourcesListSqlContainers(resp *http.Response) (result SqlResourcesListSqlContainersOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
