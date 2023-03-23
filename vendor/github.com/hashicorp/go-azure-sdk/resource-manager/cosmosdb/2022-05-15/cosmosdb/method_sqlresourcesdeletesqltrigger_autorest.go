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

type SqlResourcesDeleteSqlTriggerOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesDeleteSqlTrigger ...
func (c CosmosDBClient) SqlResourcesDeleteSqlTrigger(ctx context.Context, id TriggerId) (result SqlResourcesDeleteSqlTriggerOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesDeleteSqlTrigger(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlTrigger", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesDeleteSqlTrigger(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesDeleteSqlTrigger", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesDeleteSqlTriggerThenPoll performs SqlResourcesDeleteSqlTrigger then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlTriggerThenPoll(ctx context.Context, id TriggerId) error {
	result, err := c.SqlResourcesDeleteSqlTrigger(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlTrigger: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlTrigger: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesDeleteSqlTrigger prepares the SqlResourcesDeleteSqlTrigger request.
func (c CosmosDBClient) preparerForSqlResourcesDeleteSqlTrigger(ctx context.Context, id TriggerId) (*http.Request, error) {
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

// senderForSqlResourcesDeleteSqlTrigger sends the SqlResourcesDeleteSqlTrigger request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesDeleteSqlTrigger(ctx context.Context, req *http.Request) (future SqlResourcesDeleteSqlTriggerOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
