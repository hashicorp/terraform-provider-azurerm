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

type DatabaseAccountsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsCreateOrUpdate ...
func (c CosmosDBClient) DatabaseAccountsCreateOrUpdate(ctx context.Context, id DatabaseAccountId, input DatabaseAccountCreateUpdateParameters) (result DatabaseAccountsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsCreateOrUpdateThenPoll performs DatabaseAccountsCreateOrUpdate then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsCreateOrUpdateThenPoll(ctx context.Context, id DatabaseAccountId, input DatabaseAccountCreateUpdateParameters) error {
	result, err := c.DatabaseAccountsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsCreateOrUpdate prepares the DatabaseAccountsCreateOrUpdate request.
func (c CosmosDBClient) preparerForDatabaseAccountsCreateOrUpdate(ctx context.Context, id DatabaseAccountId, input DatabaseAccountCreateUpdateParameters) (*http.Request, error) {
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

// senderForDatabaseAccountsCreateOrUpdate sends the DatabaseAccountsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsCreateOrUpdate(ctx context.Context, req *http.Request) (future DatabaseAccountsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
