package diskaccesses

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

type DeleteAPrivateEndpointConnectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeleteAPrivateEndpointConnection ...
func (c DiskAccessesClient) DeleteAPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId) (result DeleteAPrivateEndpointConnectionOperationResponse, err error) {
	req, err := c.preparerForDeleteAPrivateEndpointConnection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "DeleteAPrivateEndpointConnection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeleteAPrivateEndpointConnection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "DeleteAPrivateEndpointConnection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeleteAPrivateEndpointConnectionThenPoll performs DeleteAPrivateEndpointConnection then polls until it's completed
func (c DiskAccessesClient) DeleteAPrivateEndpointConnectionThenPoll(ctx context.Context, id PrivateEndpointConnectionId) error {
	result, err := c.DeleteAPrivateEndpointConnection(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeleteAPrivateEndpointConnection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeleteAPrivateEndpointConnection: %+v", err)
	}

	return nil
}

// preparerForDeleteAPrivateEndpointConnection prepares the DeleteAPrivateEndpointConnection request.
func (c DiskAccessesClient) preparerForDeleteAPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId) (*http.Request, error) {
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

// senderForDeleteAPrivateEndpointConnection sends the DeleteAPrivateEndpointConnection request. The method will close the
// http.Response Body if it receives an error.
func (c DiskAccessesClient) senderForDeleteAPrivateEndpointConnection(ctx context.Context, req *http.Request) (future DeleteAPrivateEndpointConnectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
