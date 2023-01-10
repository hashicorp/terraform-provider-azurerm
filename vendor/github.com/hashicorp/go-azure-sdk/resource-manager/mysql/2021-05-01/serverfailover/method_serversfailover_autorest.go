package serverfailover

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

type ServersFailoverOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServersFailover ...
func (c ServerFailoverClient) ServersFailover(ctx context.Context, id FlexibleServerId) (result ServersFailoverOperationResponse, err error) {
	req, err := c.preparerForServersFailover(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverfailover.ServerFailoverClient", "ServersFailover", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServersFailover(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serverfailover.ServerFailoverClient", "ServersFailover", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServersFailoverThenPoll performs ServersFailover then polls until it's completed
func (c ServerFailoverClient) ServersFailoverThenPoll(ctx context.Context, id FlexibleServerId) error {
	result, err := c.ServersFailover(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ServersFailover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServersFailover: %+v", err)
	}

	return nil
}

// preparerForServersFailover prepares the ServersFailover request.
func (c ServerFailoverClient) preparerForServersFailover(ctx context.Context, id FlexibleServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/failover", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServersFailover sends the ServersFailover request. The method will close the
// http.Response Body if it receives an error.
func (c ServerFailoverClient) senderForServersFailover(ctx context.Context, req *http.Request) (future ServersFailoverOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
