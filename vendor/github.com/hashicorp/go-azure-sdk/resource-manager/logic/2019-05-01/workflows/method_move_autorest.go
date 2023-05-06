package workflows

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

type MoveOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Move ...
func (c WorkflowsClient) Move(ctx context.Context, id WorkflowId, input WorkflowReference) (result MoveOperationResponse, err error) {
	req, err := c.preparerForMove(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "Move", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMove(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "Move", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MoveThenPoll performs Move then polls until it's completed
func (c WorkflowsClient) MoveThenPoll(ctx context.Context, id WorkflowId, input WorkflowReference) error {
	result, err := c.Move(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Move: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Move: %+v", err)
	}

	return nil
}

// preparerForMove prepares the Move request.
func (c WorkflowsClient) preparerForMove(ctx context.Context, id WorkflowId, input WorkflowReference) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/move", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMove sends the Move request. The method will close the
// http.Response Body if it receives an error.
func (c WorkflowsClient) senderForMove(ctx context.Context, req *http.Request) (future MoveOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
