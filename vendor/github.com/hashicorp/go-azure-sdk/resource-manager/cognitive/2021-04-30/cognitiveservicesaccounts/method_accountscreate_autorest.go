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

type AccountsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AccountsCreate ...
func (c CognitiveServicesAccountsClient) AccountsCreate(ctx context.Context, id AccountId, input Account) (result AccountsCreateOperationResponse, err error) {
	req, err := c.preparerForAccountsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAccountsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AccountsCreateThenPoll performs AccountsCreate then polls until it's completed
func (c CognitiveServicesAccountsClient) AccountsCreateThenPoll(ctx context.Context, id AccountId, input Account) error {
	result, err := c.AccountsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AccountsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AccountsCreate: %+v", err)
	}

	return nil
}

// preparerForAccountsCreate prepares the AccountsCreate request.
func (c CognitiveServicesAccountsClient) preparerForAccountsCreate(ctx context.Context, id AccountId, input Account) (*http.Request, error) {
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

// senderForAccountsCreate sends the AccountsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c CognitiveServicesAccountsClient) senderForAccountsCreate(ctx context.Context, req *http.Request) (future AccountsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
