package serverrestart

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

type ServersRestartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServersRestart ...
func (c ServerRestartClient) ServersRestart(ctx context.Context, id FlexibleServerId, input ServerRestartParameter) (result ServersRestartOperationResponse, err error) {
	req, err := c.preparerForServersRestart(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverrestart.ServerRestartClient", "ServersRestart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServersRestart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverrestart.ServerRestartClient", "ServersRestart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServersRestartThenPoll performs ServersRestart then polls until it's completed
func (c ServerRestartClient) ServersRestartThenPoll(ctx context.Context, id FlexibleServerId, input ServerRestartParameter) error {
	result, err := c.ServersRestart(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ServersRestart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServersRestart: %+v", err)
	}

	return nil
}

// preparerForServersRestart prepares the ServersRestart request.
func (c ServerRestartClient) preparerForServersRestart(ctx context.Context, id FlexibleServerId, input ServerRestartParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restart", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServersRestart sends the ServersRestart request. The method will close the
// http.Response Body if it receives an error.
func (c ServerRestartClient) senderForServersRestart(ctx context.Context, req *http.Request) (future ServersRestartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
