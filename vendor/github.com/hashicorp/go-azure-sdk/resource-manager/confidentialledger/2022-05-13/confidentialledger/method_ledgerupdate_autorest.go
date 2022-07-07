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

type LedgerUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LedgerUpdate ...
func (c ConfidentialLedgerClient) LedgerUpdate(ctx context.Context, id LedgerId, input ConfidentialLedger) (result LedgerUpdateOperationResponse, err error) {
	req, err := c.preparerForLedgerUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLedgerUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LedgerUpdateThenPoll performs LedgerUpdate then polls until it's completed
func (c ConfidentialLedgerClient) LedgerUpdateThenPoll(ctx context.Context, id LedgerId, input ConfidentialLedger) error {
	result, err := c.LedgerUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LedgerUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LedgerUpdate: %+v", err)
	}

	return nil
}

// preparerForLedgerUpdate prepares the LedgerUpdate request.
func (c ConfidentialLedgerClient) preparerForLedgerUpdate(ctx context.Context, id LedgerId, input ConfidentialLedger) (*http.Request, error) {
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

// senderForLedgerUpdate sends the LedgerUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ConfidentialLedgerClient) senderForLedgerUpdate(ctx context.Context, req *http.Request) (future LedgerUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
