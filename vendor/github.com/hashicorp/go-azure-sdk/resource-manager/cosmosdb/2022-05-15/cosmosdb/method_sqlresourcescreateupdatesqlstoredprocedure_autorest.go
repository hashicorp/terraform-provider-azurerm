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

type SqlResourcesCreateUpdateSqlStoredProcedureOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateSqlStoredProcedure ...
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlStoredProcedure(ctx context.Context, id StoredProcedureId, input SqlStoredProcedureCreateUpdateParameters) (result SqlResourcesCreateUpdateSqlStoredProcedureOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateSqlStoredProcedure(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlStoredProcedure", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateSqlStoredProcedure(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlStoredProcedure", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateSqlStoredProcedureThenPoll performs SqlResourcesCreateUpdateSqlStoredProcedure then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx context.Context, id StoredProcedureId, input SqlStoredProcedureCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateSqlStoredProcedure(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateSqlStoredProcedure: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateSqlStoredProcedure: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateSqlStoredProcedure prepares the SqlResourcesCreateUpdateSqlStoredProcedure request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateSqlStoredProcedure(ctx context.Context, id StoredProcedureId, input SqlStoredProcedureCreateUpdateParameters) (*http.Request, error) {
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

// senderForSqlResourcesCreateUpdateSqlStoredProcedure sends the SqlResourcesCreateUpdateSqlStoredProcedure request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateSqlStoredProcedure(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateSqlStoredProcedureOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
