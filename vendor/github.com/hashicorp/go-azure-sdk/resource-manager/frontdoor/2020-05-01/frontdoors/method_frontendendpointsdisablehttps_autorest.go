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

type FrontendEndpointsDisableHTTPSOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FrontendEndpointsDisableHTTPS ...
func (c FrontDoorsClient) FrontendEndpointsDisableHTTPS(ctx context.Context, id FrontendEndpointId) (result FrontendEndpointsDisableHTTPSOperationResponse, err error) {
	req, err := c.preparerForFrontendEndpointsDisableHTTPS(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsDisableHTTPS", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFrontendEndpointsDisableHTTPS(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsDisableHTTPS", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FrontendEndpointsDisableHTTPSThenPoll performs FrontendEndpointsDisableHTTPS then polls until it's completed
func (c FrontDoorsClient) FrontendEndpointsDisableHTTPSThenPoll(ctx context.Context, id FrontendEndpointId) error {
	result, err := c.FrontendEndpointsDisableHTTPS(ctx, id)
	if err != nil {
		return fmt.Errorf("performing FrontendEndpointsDisableHTTPS: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FrontendEndpointsDisableHTTPS: %+v", err)
	}

	return nil
}

// preparerForFrontendEndpointsDisableHTTPS prepares the FrontendEndpointsDisableHTTPS request.
func (c FrontDoorsClient) preparerForFrontendEndpointsDisableHTTPS(ctx context.Context, id FrontendEndpointId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disableHttps", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFrontendEndpointsDisableHTTPS sends the FrontendEndpointsDisableHTTPS request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForFrontendEndpointsDisableHTTPS(ctx context.Context, req *http.Request) (future FrontendEndpointsDisableHTTPSOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
