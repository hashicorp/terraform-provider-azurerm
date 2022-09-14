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

type DatabaseAccountsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsUpdate ...
func (c CosmosDBClient) DatabaseAccountsUpdate(ctx context.Context, id DatabaseAccountId, input DatabaseAccountUpdateParameters) (result DatabaseAccountsUpdateOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsUpdateThenPoll performs DatabaseAccountsUpdate then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsUpdateThenPoll(ctx context.Context, id DatabaseAccountId, input DatabaseAccountUpdateParameters) error {
	result, err := c.DatabaseAccountsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsUpdate: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsUpdate prepares the DatabaseAccountsUpdate request.
func (c CosmosDBClient) preparerForDatabaseAccountsUpdate(ctx context.Context, id DatabaseAccountId, input DatabaseAccountUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabaseAccountsUpdate sends the DatabaseAccountsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsUpdate(ctx context.Context, req *http.Request) (future DatabaseAccountsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
