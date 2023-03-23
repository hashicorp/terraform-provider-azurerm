package nodetype

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

type RestartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Restart ...
func (c NodeTypeClient) Restart(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) (result RestartOperationResponse, err error) {
	req, err := c.preparerForRestart(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nodetype.NodeTypeClient", "Restart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRestart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nodetype.NodeTypeClient", "Restart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RestartThenPoll performs Restart then polls until it's completed
func (c NodeTypeClient) RestartThenPoll(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) error {
	result, err := c.Restart(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Restart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Restart: %+v", err)
	}

	return nil
}

// preparerForRestart prepares the Restart request.
func (c NodeTypeClient) preparerForRestart(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) (*http.Request, error) {
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

// senderForRestart sends the Restart request. The method will close the
// http.Response Body if it receives an error.
func (c NodeTypeClient) senderForRestart(ctx context.Context, req *http.Request) (future RestartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
