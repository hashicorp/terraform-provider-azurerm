package liveevents

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

type ResetOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Reset ...
func (c LiveEventsClient) Reset(ctx context.Context, id LiveEventId) (result ResetOperationResponse, err error) {
	req, err := c.preparerForReset(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Reset", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReset(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Reset", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResetThenPoll performs Reset then polls until it's completed
func (c LiveEventsClient) ResetThenPoll(ctx context.Context, id LiveEventId) error {
	result, err := c.Reset(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Reset: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Reset: %+v", err)
	}

	return nil
}

// preparerForReset prepares the Reset request.
func (c LiveEventsClient) preparerForReset(ctx context.Context, id LiveEventId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reset", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReset sends the Reset request. The method will close the
// http.Response Body if it receives an error.
func (c LiveEventsClient) senderForReset(ctx context.Context, req *http.Request) (future ResetOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
