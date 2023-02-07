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

type DigitalTwinsEndpointCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DigitalTwinsEndpointCreateOrUpdate ...
func (c EndpointsClient) DigitalTwinsEndpointCreateOrUpdate(ctx context.Context, id EndpointId, input DigitalTwinsEndpointResource) (result DigitalTwinsEndpointCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsEndpointCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDigitalTwinsEndpointCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DigitalTwinsEndpointCreateOrUpdateThenPoll performs DigitalTwinsEndpointCreateOrUpdate then polls until it's completed
func (c EndpointsClient) DigitalTwinsEndpointCreateOrUpdateThenPoll(ctx context.Context, id EndpointId, input DigitalTwinsEndpointResource) error {
	result, err := c.DigitalTwinsEndpointCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsEndpointCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DigitalTwinsEndpointCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForDigitalTwinsEndpointCreateOrUpdate prepares the DigitalTwinsEndpointCreateOrUpdate request.
func (c EndpointsClient) preparerForDigitalTwinsEndpointCreateOrUpdate(ctx context.Context, id EndpointId, input DigitalTwinsEndpointResource) (*http.Request, error) {
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

// senderForDigitalTwinsEndpointCreateOrUpdate sends the DigitalTwinsEndpointCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EndpointsClient) senderForDigitalTwinsEndpointCreateOrUpdate(ctx context.Context, req *http.Request) (future DigitalTwinsEndpointCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
