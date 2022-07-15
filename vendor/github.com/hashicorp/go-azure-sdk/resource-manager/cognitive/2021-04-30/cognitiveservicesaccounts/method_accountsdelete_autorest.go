package cognitiveservicesaccounts

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

type AccountsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AccountsDelete ...
func (c CognitiveServicesAccountsClient) AccountsDelete(ctx context.Context, id AccountId) (result AccountsDeleteOperationResponse, err error) {
	req, err := c.preparerForAccountsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAccountsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AccountsDeleteThenPoll performs AccountsDelete then polls until it's completed
func (c CognitiveServicesAccountsClient) AccountsDeleteThenPoll(ctx context.Context, id AccountId) error {
	result, err := c.AccountsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AccountsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AccountsDelete: %+v", err)
	}

	return nil
}

// preparerForAccountsDelete prepares the AccountsDelete request.
func (c CognitiveServicesAccountsClient) preparerForAccountsDelete(ctx context.Context, id AccountId) (*http.Request, error) {
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

// senderForAccountsDelete sends the AccountsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c CognitiveServicesAccountsClient) senderForAccountsDelete(ctx context.Context, req *http.Request) (future AccountsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
