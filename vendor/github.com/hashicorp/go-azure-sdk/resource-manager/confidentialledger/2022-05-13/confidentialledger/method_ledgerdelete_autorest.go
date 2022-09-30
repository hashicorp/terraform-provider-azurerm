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

type LedgerDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LedgerDelete ...
func (c ConfidentialLedgerClient) LedgerDelete(ctx context.Context, id LedgerId) (result LedgerDeleteOperationResponse, err error) {
	req, err := c.preparerForLedgerDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLedgerDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LedgerDeleteThenPoll performs LedgerDelete then polls until it's completed
func (c ConfidentialLedgerClient) LedgerDeleteThenPoll(ctx context.Context, id LedgerId) error {
	result, err := c.LedgerDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing LedgerDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LedgerDelete: %+v", err)
	}

	return nil
}

// preparerForLedgerDelete prepares the LedgerDelete request.
func (c ConfidentialLedgerClient) preparerForLedgerDelete(ctx context.Context, id LedgerId) (*http.Request, error) {
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

// senderForLedgerDelete sends the LedgerDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ConfidentialLedgerClient) senderForLedgerDelete(ctx context.Context, req *http.Request) (future LedgerDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
