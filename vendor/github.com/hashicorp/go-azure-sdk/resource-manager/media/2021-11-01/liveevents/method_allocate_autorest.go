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

type AllocateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Allocate ...
func (c LiveEventsClient) Allocate(ctx context.Context, id LiveEventId) (result AllocateOperationResponse, err error) {
	req, err := c.preparerForAllocate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Allocate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAllocate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Allocate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AllocateThenPoll performs Allocate then polls until it's completed
func (c LiveEventsClient) AllocateThenPoll(ctx context.Context, id LiveEventId) error {
	result, err := c.Allocate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Allocate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Allocate: %+v", err)
	}

	return nil
}

// preparerForAllocate prepares the Allocate request.
func (c LiveEventsClient) preparerForAllocate(ctx context.Context, id LiveEventId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/allocate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAllocate sends the Allocate request. The method will close the
// http.Response Body if it receives an error.
func (c LiveEventsClient) senderForAllocate(ctx context.Context, req *http.Request) (future AllocateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
