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

type SqlResourcesDeleteSqlDatabaseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesDeleteSqlDatabase ...
func (c CosmosDBClient) SqlResourcesDeleteSqlDatabase(ctx context.Context, id SqlDatabaseId) (result SqlResourcesDeleteSqlDatabaseOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesDeleteSqlDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlDatabase", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesDeleteSqlDatabase(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesDeleteSqlDatabaseThenPoll performs SqlResourcesDeleteSqlDatabase then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlDatabaseThenPoll(ctx context.Context, id SqlDatabaseId) error {
	result, err := c.SqlResourcesDeleteSqlDatabase(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlDatabase: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesDeleteSqlDatabase prepares the SqlResourcesDeleteSqlDatabase request.
func (c CosmosDBClient) preparerForSqlResourcesDeleteSqlDatabase(ctx context.Context, id SqlDatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSqlResourcesDeleteSqlDatabase sends the SqlResourcesDeleteSqlDatabase request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesDeleteSqlDatabase(ctx context.Context, req *http.Request) (future SqlResourcesDeleteSqlDatabaseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
