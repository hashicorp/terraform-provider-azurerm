package frontdoors

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

type FrontendEndpointsEnableHTTPSOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FrontendEndpointsEnableHTTPS ...
func (c FrontDoorsClient) FrontendEndpointsEnableHTTPS(ctx context.Context, id FrontendEndpointId, input CustomHTTPSConfiguration) (result FrontendEndpointsEnableHTTPSOperationResponse, err error) {
	req, err := c.preparerForFrontendEndpointsEnableHTTPS(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsEnableHTTPS", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFrontendEndpointsEnableHTTPS(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsEnableHTTPS", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FrontendEndpointsEnableHTTPSThenPoll performs FrontendEndpointsEnableHTTPS then polls until it's completed
func (c FrontDoorsClient) FrontendEndpointsEnableHTTPSThenPoll(ctx context.Context, id FrontendEndpointId, input CustomHTTPSConfiguration) error {
	result, err := c.FrontendEndpointsEnableHTTPS(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing FrontendEndpointsEnableHTTPS: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FrontendEndpointsEnableHTTPS: %+v", err)
	}

	return nil
}

// preparerForFrontendEndpointsEnableHTTPS prepares the FrontendEndpointsEnableHTTPS request.
func (c FrontDoorsClient) preparerForFrontendEndpointsEnableHTTPS(ctx context.Context, id FrontendEndpointId, input CustomHTTPSConfiguration) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enableHttps", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFrontendEndpointsEnableHTTPS sends the FrontendEndpointsEnableHTTPS request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForFrontendEndpointsEnableHTTPS(ctx context.Context, req *http.Request) (future FrontendEndpointsEnableHTTPSOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
