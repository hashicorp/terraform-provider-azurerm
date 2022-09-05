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

type DatabaseAccountsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsDelete ...
func (c CosmosDBClient) DatabaseAccountsDelete(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsDeleteOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsDeleteThenPoll performs DatabaseAccountsDelete then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsDeleteThenPoll(ctx context.Context, id DatabaseAccountId) error {
	result, err := c.DatabaseAccountsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsDelete: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsDelete prepares the DatabaseAccountsDelete request.
func (c CosmosDBClient) preparerForDatabaseAccountsDelete(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
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

// senderForDatabaseAccountsDelete sends the DatabaseAccountsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsDelete(ctx context.Context, req *http.Request) (future DatabaseAccountsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
