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

type UpdateAPrivateEndpointConnectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdateAPrivateEndpointConnection ...
func (c DiskAccessesClient) UpdateAPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (result UpdateAPrivateEndpointConnectionOperationResponse, err error) {
	req, err := c.preparerForUpdateAPrivateEndpointConnection(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "UpdateAPrivateEndpointConnection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdateAPrivateEndpointConnection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "UpdateAPrivateEndpointConnection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdateAPrivateEndpointConnectionThenPoll performs UpdateAPrivateEndpointConnection then polls until it's completed
func (c DiskAccessesClient) UpdateAPrivateEndpointConnectionThenPoll(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) error {
	result, err := c.UpdateAPrivateEndpointConnection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UpdateAPrivateEndpointConnection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdateAPrivateEndpointConnection: %+v", err)
	}

	return nil
}

// preparerForUpdateAPrivateEndpointConnection prepares the UpdateAPrivateEndpointConnection request.
func (c DiskAccessesClient) preparerForUpdateAPrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (*http.Request, error) {
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

// senderForUpdateAPrivateEndpointConnection sends the UpdateAPrivateEndpointConnection request. The method will close the
// http.Response Body if it receives an error.
func (c DiskAccessesClient) senderForUpdateAPrivateEndpointConnection(ctx context.Context, req *http.Request) (future UpdateAPrivateEndpointConnectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
