package protectioncontainers

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

type UnregisterOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Unregister ...
func (c ProtectionContainersClient) Unregister(ctx context.Context, id ProtectionContainerId) (result UnregisterOperationResponse, err error) {
	req, err := c.preparerForUnregister(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Unregister", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUnregister(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Unregister", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UnregisterThenPoll performs Unregister then polls until it's completed
func (c ProtectionContainersClient) UnregisterThenPoll(ctx context.Context, id ProtectionContainerId) error {
	return c.UnregisterCallbackThenPoll(ctx, id, nil)
}

// UnregisterCallbackThenPoll performs Unregister, runs the optional callback function, then polls until it's completed
func (c ProtectionContainersClient) UnregisterCallbackThenPoll(ctx context.Context, id ProtectionContainerId, callback func() error) error {
	result, err := c.Unregister(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Unregister: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Unregister: %+v", err)
	}

	return nil
}

// preparerForUnregister prepares the Unregister request.
func (c ProtectionContainersClient) preparerForUnregister(ctx context.Context, id ProtectionContainerId) (*http.Request, error) {
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

// senderForUnregister sends the Unregister request. The method will close the
// http.Response Body if it receives an error.
func (c ProtectionContainersClient) senderForUnregister(ctx context.Context, req *http.Request) (future UnregisterOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
