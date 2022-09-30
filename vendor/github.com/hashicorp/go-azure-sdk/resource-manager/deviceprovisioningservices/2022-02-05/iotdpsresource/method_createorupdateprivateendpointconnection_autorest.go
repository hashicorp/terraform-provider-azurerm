package iotdpsresource

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

type CreateOrUpdatePrivateEndpointConnectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CreateOrUpdatePrivateEndpointConnection ...
func (c IotDpsResourceClient) CreateOrUpdatePrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (result CreateOrUpdatePrivateEndpointConnectionOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdatePrivateEndpointConnection(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "CreateOrUpdatePrivateEndpointConnection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateOrUpdatePrivateEndpointConnection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "CreateOrUpdatePrivateEndpointConnection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePrivateEndpointConnectionThenPoll performs CreateOrUpdatePrivateEndpointConnection then polls until it's completed
func (c IotDpsResourceClient) CreateOrUpdatePrivateEndpointConnectionThenPoll(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) error {
	result, err := c.CreateOrUpdatePrivateEndpointConnection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdatePrivateEndpointConnection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateOrUpdatePrivateEndpointConnection: %+v", err)
	}

	return nil
}

// preparerForCreateOrUpdatePrivateEndpointConnection prepares the CreateOrUpdatePrivateEndpointConnection request.
func (c IotDpsResourceClient) preparerForCreateOrUpdatePrivateEndpointConnection(ctx context.Context, id PrivateEndpointConnectionId, input PrivateEndpointConnection) (*http.Request, error) {
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

// senderForCreateOrUpdatePrivateEndpointConnection sends the CreateOrUpdatePrivateEndpointConnection request. The method will close the
// http.Response Body if it receives an error.
func (c IotDpsResourceClient) senderForCreateOrUpdatePrivateEndpointConnection(ctx context.Context, req *http.Request) (future CreateOrUpdatePrivateEndpointConnectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
