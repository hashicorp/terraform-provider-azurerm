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

type SqlResourcesCreateUpdateSqlTriggerOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateSqlTrigger ...
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlTrigger(ctx context.Context, id TriggerId, input SqlTriggerCreateUpdateParameters) (result SqlResourcesCreateUpdateSqlTriggerOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateSqlTrigger(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlTrigger", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateSqlTrigger(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateSqlTrigger", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateSqlTriggerThenPoll performs SqlResourcesCreateUpdateSqlTrigger then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateSqlTriggerThenPoll(ctx context.Context, id TriggerId, input SqlTriggerCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateSqlTrigger(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateSqlTrigger: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateSqlTrigger: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateSqlTrigger prepares the SqlResourcesCreateUpdateSqlTrigger request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateSqlTrigger(ctx context.Context, id TriggerId, input SqlTriggerCreateUpdateParameters) (*http.Request, error) {
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

// senderForSqlResourcesCreateUpdateSqlTrigger sends the SqlResourcesCreateUpdateSqlTrigger request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateSqlTrigger(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateSqlTriggerOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
