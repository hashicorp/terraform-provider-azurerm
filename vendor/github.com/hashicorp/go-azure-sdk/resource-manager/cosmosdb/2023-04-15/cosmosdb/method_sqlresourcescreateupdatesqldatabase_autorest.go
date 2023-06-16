package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesCreateUpdateSqlDatabaseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateSqlDatabase ...
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlDatabase(ctx context.Context, id SqlDatabaseId, input SqlDatabaseCreateUpdateParameters) (result SqlResourcesCreateUpdateSqlDatabaseOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateSqlDatabase(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlDatabase", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateSqlDatabase(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateSqlDatabaseThenPoll performs SqlResourcesCreateUpdateSqlDatabase then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlDatabaseThenPoll(ctx context.Context, id SqlDatabaseId, input SqlDatabaseCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateSqlDatabase(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateSqlDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateSqlDatabase: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateSqlDatabase prepares the SqlResourcesCreateUpdateSqlDatabase request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateSqlDatabase(ctx context.Context, id SqlDatabaseId, input SqlDatabaseCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSqlResourcesCreateUpdateSqlDatabase sends the SqlResourcesCreateUpdateSqlDatabase request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateSqlDatabase(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateSqlDatabaseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
