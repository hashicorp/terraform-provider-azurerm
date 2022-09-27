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

type SqlResourcesDeleteSqlStoredProcedureOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesDeleteSqlStoredProcedure ...
func (c CosmosDBClient) SqlResourcesDeleteSqlStoredProcedure(ctx context.Context, id StoredProcedureId) (result SqlResourcesDeleteSqlStoredProcedureOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesDeleteSqlStoredProcedure(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlStoredProcedure", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesDeleteSqlStoredProcedure(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlStoredProcedure", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesDeleteSqlStoredProcedureThenPoll performs SqlResourcesDeleteSqlStoredProcedure then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlStoredProcedureThenPoll(ctx context.Context, id StoredProcedureId) error {
	result, err := c.SqlResourcesDeleteSqlStoredProcedure(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlStoredProcedure: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlStoredProcedure: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesDeleteSqlStoredProcedure prepares the SqlResourcesDeleteSqlStoredProcedure request.
func (c CosmosDBClient) preparerForSqlResourcesDeleteSqlStoredProcedure(ctx context.Context, id StoredProcedureId) (*http.Request, error) {
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

// senderForSqlResourcesDeleteSqlStoredProcedure sends the SqlResourcesDeleteSqlStoredProcedure request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesDeleteSqlStoredProcedure(ctx context.Context, req *http.Request) (future SqlResourcesDeleteSqlStoredProcedureOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
