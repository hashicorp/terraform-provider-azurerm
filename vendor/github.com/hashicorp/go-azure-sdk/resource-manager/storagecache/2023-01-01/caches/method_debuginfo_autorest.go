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

type DebugInfoOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DebugInfo ...
func (c CachesClient) DebugInfo(ctx context.Context, id CacheId) (result DebugInfoOperationResponse, err error) {
	req, err := c.preparerForDebugInfo(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "DebugInfo", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDebugInfo(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "DebugInfo", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DebugInfoThenPoll performs DebugInfo then polls until it's completed
func (c CachesClient) DebugInfoThenPoll(ctx context.Context, id CacheId) error {
	result, err := c.DebugInfo(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DebugInfo: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DebugInfo: %+v", err)
	}

	return nil
}

// preparerForDebugInfo prepares the DebugInfo request.
func (c CachesClient) preparerForDebugInfo(ctx context.Context, id CacheId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/debugInfo", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDebugInfo sends the DebugInfo request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForDebugInfo(ctx context.Context, req *http.Request) (future DebugInfoOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
