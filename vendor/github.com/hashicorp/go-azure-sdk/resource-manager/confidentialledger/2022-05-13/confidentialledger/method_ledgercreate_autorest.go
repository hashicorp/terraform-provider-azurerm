package confidentialledger

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

type LedgerCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LedgerCreate ...
func (c ConfidentialLedgerClient) LedgerCreate(ctx context.Context, id LedgerId, input ConfidentialLedger) (result LedgerCreateOperationResponse, err error) {
	req, err := c.preparerForLedgerCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLedgerCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LedgerCreateThenPoll performs LedgerCreate then polls until it's completed
func (c ConfidentialLedgerClient) LedgerCreateThenPoll(ctx context.Context, id LedgerId, input ConfidentialLedger) error {
	result, err := c.LedgerCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LedgerCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LedgerCreate: %+v", err)
	}

	return nil
}

// preparerForLedgerCreate prepares the LedgerCreate request.
func (c ConfidentialLedgerClient) preparerForLedgerCreate(ctx context.Context, id LedgerId, input ConfidentialLedger) (*http.Request, error) {
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

// senderForLedgerCreate sends the LedgerCreate request. The method will close the
// http.Response Body if it receives an error.
func (c ConfidentialLedgerClient) senderForLedgerCreate(ctx context.Context, req *http.Request) (future LedgerCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
