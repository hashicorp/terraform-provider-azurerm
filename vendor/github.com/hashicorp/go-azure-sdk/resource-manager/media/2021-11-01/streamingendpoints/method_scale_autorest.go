package streamingendpoints

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

type ScaleOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Scale ...
func (c StreamingEndpointsClient) Scale(ctx context.Context, id StreamingEndpointId, input StreamingEntityScaleUnit) (result ScaleOperationResponse, err error) {
	req, err := c.preparerForScale(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "Scale", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForScale(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "Scale", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ScaleThenPoll performs Scale then polls until it's completed
func (c StreamingEndpointsClient) ScaleThenPoll(ctx context.Context, id StreamingEndpointId, input StreamingEntityScaleUnit) error {
	result, err := c.Scale(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Scale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Scale: %+v", err)
	}

	return nil
}

// preparerForScale prepares the Scale request.
func (c StreamingEndpointsClient) preparerForScale(ctx context.Context, id StreamingEndpointId, input StreamingEntityScaleUnit) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/scale", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForScale sends the Scale request. The method will close the
// http.Response Body if it receives an error.
func (c StreamingEndpointsClient) senderForScale(ctx context.Context, req *http.Request) (future ScaleOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
