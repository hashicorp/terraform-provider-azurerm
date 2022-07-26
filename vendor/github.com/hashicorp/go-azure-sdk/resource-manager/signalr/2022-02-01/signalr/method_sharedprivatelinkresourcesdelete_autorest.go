package signalr

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

type SharedPrivateLinkResourcesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SharedPrivateLinkResourcesDelete ...
func (c SignalRClient) SharedPrivateLinkResourcesDelete(ctx context.Context, id SharedPrivateLinkResourceId) (result SharedPrivateLinkResourcesDeleteOperationResponse, err error) {
	req, err := c.preparerForSharedPrivateLinkResourcesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSharedPrivateLinkResourcesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SharedPrivateLinkResourcesDeleteThenPoll performs SharedPrivateLinkResourcesDelete then polls until it's completed
func (c SignalRClient) SharedPrivateLinkResourcesDeleteThenPoll(ctx context.Context, id SharedPrivateLinkResourceId) error {
	result, err := c.SharedPrivateLinkResourcesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SharedPrivateLinkResourcesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SharedPrivateLinkResourcesDelete: %+v", err)
	}

	return nil
}

// preparerForSharedPrivateLinkResourcesDelete prepares the SharedPrivateLinkResourcesDelete request.
func (c SignalRClient) preparerForSharedPrivateLinkResourcesDelete(ctx context.Context, id SharedPrivateLinkResourceId) (*http.Request, error) {
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

// senderForSharedPrivateLinkResourcesDelete sends the SharedPrivateLinkResourcesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForSharedPrivateLinkResourcesDelete(ctx context.Context, req *http.Request) (future SharedPrivateLinkResourcesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
