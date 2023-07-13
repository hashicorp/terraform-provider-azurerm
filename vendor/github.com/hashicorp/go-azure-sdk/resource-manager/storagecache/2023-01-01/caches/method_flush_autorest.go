package caches

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

type FlushOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Flush ...
func (c CachesClient) Flush(ctx context.Context, id CacheId) (result FlushOperationResponse, err error) {
	req, err := c.preparerForFlush(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "Flush", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFlush(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "Flush", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FlushThenPoll performs Flush then polls until it's completed
func (c CachesClient) FlushThenPoll(ctx context.Context, id CacheId) error {
	result, err := c.Flush(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Flush: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Flush: %+v", err)
	}

	return nil
}

// preparerForFlush prepares the Flush request.
func (c CachesClient) preparerForFlush(ctx context.Context, id CacheId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/flush", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFlush sends the Flush request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForFlush(ctx context.Context, req *http.Request) (future FlushOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
