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

type SharedPrivateLinkResourcesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SharedPrivateLinkResourcesCreateOrUpdate ...
func (c SignalRClient) SharedPrivateLinkResourcesCreateOrUpdate(ctx context.Context, id SharedPrivateLinkResourceId, input SharedPrivateLinkResource) (result SharedPrivateLinkResourcesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForSharedPrivateLinkResourcesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSharedPrivateLinkResourcesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "SharedPrivateLinkResourcesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SharedPrivateLinkResourcesCreateOrUpdateThenPoll performs SharedPrivateLinkResourcesCreateOrUpdate then polls until it's completed
func (c SignalRClient) SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx context.Context, id SharedPrivateLinkResourceId, input SharedPrivateLinkResource) error {
	result, err := c.SharedPrivateLinkResourcesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SharedPrivateLinkResourcesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SharedPrivateLinkResourcesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForSharedPrivateLinkResourcesCreateOrUpdate prepares the SharedPrivateLinkResourcesCreateOrUpdate request.
func (c SignalRClient) preparerForSharedPrivateLinkResourcesCreateOrUpdate(ctx context.Context, id SharedPrivateLinkResourceId, input SharedPrivateLinkResource) (*http.Request, error) {
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

// senderForSharedPrivateLinkResourcesCreateOrUpdate sends the SharedPrivateLinkResourcesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForSharedPrivateLinkResourcesCreateOrUpdate(ctx context.Context, req *http.Request) (future SharedPrivateLinkResourcesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
