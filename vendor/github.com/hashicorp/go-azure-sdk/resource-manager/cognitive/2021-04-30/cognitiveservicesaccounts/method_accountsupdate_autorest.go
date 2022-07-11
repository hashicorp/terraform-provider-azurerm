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

type AccountsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AccountsUpdate ...
func (c CognitiveServicesAccountsClient) AccountsUpdate(ctx context.Context, id AccountId, input Account) (result AccountsUpdateOperationResponse, err error) {
	req, err := c.preparerForAccountsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAccountsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AccountsUpdateThenPoll performs AccountsUpdate then polls until it's completed
func (c CognitiveServicesAccountsClient) AccountsUpdateThenPoll(ctx context.Context, id AccountId, input Account) error {
	result, err := c.AccountsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AccountsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AccountsUpdate: %+v", err)
	}

	return nil
}

// preparerForAccountsUpdate prepares the AccountsUpdate request.
func (c CognitiveServicesAccountsClient) preparerForAccountsUpdate(ctx context.Context, id AccountId, input Account) (*http.Request, error) {
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

// senderForAccountsUpdate sends the AccountsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c CognitiveServicesAccountsClient) senderForAccountsUpdate(ctx context.Context, req *http.Request) (future AccountsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
