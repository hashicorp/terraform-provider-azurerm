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

type EndpointsPurgeContentOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// EndpointsPurgeContent ...
func (c FrontDoorsClient) EndpointsPurgeContent(ctx context.Context, id FrontDoorId, input PurgeParameters) (result EndpointsPurgeContentOperationResponse, err error) {
	req, err := c.preparerForEndpointsPurgeContent(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "EndpointsPurgeContent", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForEndpointsPurgeContent(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "EndpointsPurgeContent", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// EndpointsPurgeContentThenPoll performs EndpointsPurgeContent then polls until it's completed
func (c FrontDoorsClient) EndpointsPurgeContentThenPoll(ctx context.Context, id FrontDoorId, input PurgeParameters) error {
	result, err := c.EndpointsPurgeContent(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing EndpointsPurgeContent: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after EndpointsPurgeContent: %+v", err)
	}

	return nil
}

// preparerForEndpointsPurgeContent prepares the EndpointsPurgeContent request.
func (c FrontDoorsClient) preparerForEndpointsPurgeContent(ctx context.Context, id FrontDoorId, input PurgeParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/purge", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForEndpointsPurgeContent sends the EndpointsPurgeContent request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForEndpointsPurgeContent(ctx context.Context, req *http.Request) (future EndpointsPurgeContentOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
