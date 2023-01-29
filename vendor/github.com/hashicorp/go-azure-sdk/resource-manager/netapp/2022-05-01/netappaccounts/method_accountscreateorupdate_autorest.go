package netappaccounts

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

type AccountsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AccountsCreateOrUpdate ...
func (c NetAppAccountsClient) AccountsCreateOrUpdate(ctx context.Context, id NetAppAccountId, input NetAppAccount) (result AccountsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAccountsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAccountsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AccountsCreateOrUpdateThenPoll performs AccountsCreateOrUpdate then polls until it's completed
func (c NetAppAccountsClient) AccountsCreateOrUpdateThenPoll(ctx context.Context, id NetAppAccountId, input NetAppAccount) error {
	result, err := c.AccountsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AccountsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AccountsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForAccountsCreateOrUpdate prepares the AccountsCreateOrUpdate request.
func (c NetAppAccountsClient) preparerForAccountsCreateOrUpdate(ctx context.Context, id NetAppAccountId, input NetAppAccount) (*http.Request, error) {
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

// senderForAccountsCreateOrUpdate sends the AccountsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c NetAppAccountsClient) senderForAccountsCreateOrUpdate(ctx context.Context, req *http.Request) (future AccountsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
