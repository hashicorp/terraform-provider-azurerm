package endpoints

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

type DigitalTwinsEndpointDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DigitalTwinsEndpointDelete ...
func (c EndpointsClient) DigitalTwinsEndpointDelete(ctx context.Context, id EndpointId) (result DigitalTwinsEndpointDeleteOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsEndpointDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDigitalTwinsEndpointDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DigitalTwinsEndpointDeleteThenPoll performs DigitalTwinsEndpointDelete then polls until it's completed
func (c EndpointsClient) DigitalTwinsEndpointDeleteThenPoll(ctx context.Context, id EndpointId) error {
	result, err := c.DigitalTwinsEndpointDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsEndpointDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DigitalTwinsEndpointDelete: %+v", err)
	}

	return nil
}

// preparerForDigitalTwinsEndpointDelete prepares the DigitalTwinsEndpointDelete request.
func (c EndpointsClient) preparerForDigitalTwinsEndpointDelete(ctx context.Context, id EndpointId) (*http.Request, error) {
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

// senderForDigitalTwinsEndpointDelete sends the DigitalTwinsEndpointDelete request. The method will close the
// http.Response Body if it receives an error.
func (c EndpointsClient) senderForDigitalTwinsEndpointDelete(ctx context.Context, req *http.Request) (future DigitalTwinsEndpointDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
