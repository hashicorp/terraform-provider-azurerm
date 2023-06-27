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

type DatabaseAccountsRegenerateKeyOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsRegenerateKey ...
func (c CosmosDBClient) DatabaseAccountsRegenerateKey(ctx context.Context, id DatabaseAccountId, input DatabaseAccountRegenerateKeyParameters) (result DatabaseAccountsRegenerateKeyOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsRegenerateKey", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsRegenerateKey(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsRegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsRegenerateKeyThenPoll performs DatabaseAccountsRegenerateKey then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsRegenerateKeyThenPoll(ctx context.Context, id DatabaseAccountId, input DatabaseAccountRegenerateKeyParameters) error {
	result, err := c.DatabaseAccountsRegenerateKey(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsRegenerateKey: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsRegenerateKey: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsRegenerateKey prepares the DatabaseAccountsRegenerateKey request.
func (c CosmosDBClient) preparerForDatabaseAccountsRegenerateKey(ctx context.Context, id DatabaseAccountId, input DatabaseAccountRegenerateKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabaseAccountsRegenerateKey sends the DatabaseAccountsRegenerateKey request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsRegenerateKey(ctx context.Context, req *http.Request) (future DatabaseAccountsRegenerateKeyOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
