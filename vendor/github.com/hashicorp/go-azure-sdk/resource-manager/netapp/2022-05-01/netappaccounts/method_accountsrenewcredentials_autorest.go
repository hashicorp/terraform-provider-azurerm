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

type AccountsRenewCredentialsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AccountsRenewCredentials ...
func (c NetAppAccountsClient) AccountsRenewCredentials(ctx context.Context, id NetAppAccountId) (result AccountsRenewCredentialsOperationResponse, err error) {
	req, err := c.preparerForAccountsRenewCredentials(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsRenewCredentials", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAccountsRenewCredentials(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsRenewCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AccountsRenewCredentialsThenPoll performs AccountsRenewCredentials then polls until it's completed
func (c NetAppAccountsClient) AccountsRenewCredentialsThenPoll(ctx context.Context, id NetAppAccountId) error {
	result, err := c.AccountsRenewCredentials(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AccountsRenewCredentials: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AccountsRenewCredentials: %+v", err)
	}

	return nil
}

// preparerForAccountsRenewCredentials prepares the AccountsRenewCredentials request.
func (c NetAppAccountsClient) preparerForAccountsRenewCredentials(ctx context.Context, id NetAppAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/renewCredentials", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAccountsRenewCredentials sends the AccountsRenewCredentials request. The method will close the
// http.Response Body if it receives an error.
func (c NetAppAccountsClient) senderForAccountsRenewCredentials(ctx context.Context, req *http.Request) (future AccountsRenewCredentialsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
