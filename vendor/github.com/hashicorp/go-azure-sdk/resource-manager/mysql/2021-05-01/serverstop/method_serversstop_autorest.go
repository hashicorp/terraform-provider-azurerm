package serverstop

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

type ServersStopOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServersStop ...
func (c ServerStopClient) ServersStop(ctx context.Context, id FlexibleServerId) (result ServersStopOperationResponse, err error) {
	req, err := c.preparerForServersStop(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverstop.ServerStopClient", "ServersStop", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServersStop(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverstop.ServerStopClient", "ServersStop", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServersStopThenPoll performs ServersStop then polls until it's completed
func (c ServerStopClient) ServersStopThenPoll(ctx context.Context, id FlexibleServerId) error {
	result, err := c.ServersStop(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ServersStop: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServersStop: %+v", err)
	}

	return nil
}

// preparerForServersStop prepares the ServersStop request.
func (c ServerStopClient) preparerForServersStop(ctx context.Context, id FlexibleServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stop", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServersStop sends the ServersStop request. The method will close the
// http.Response Body if it receives an error.
func (c ServerStopClient) senderForServersStop(ctx context.Context, req *http.Request) (future ServersStopOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
