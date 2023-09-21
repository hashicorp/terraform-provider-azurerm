package jobs

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

type StopExecutionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StopExecution ...
func (c JobsClient) StopExecution(ctx context.Context, id ExecutionId) (result StopExecutionOperationResponse, err error) {
	req, err := c.preparerForStopExecution(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "StopExecution", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStopExecution(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "StopExecution", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StopExecutionThenPoll performs StopExecution then polls until it's completed
func (c JobsClient) StopExecutionThenPoll(ctx context.Context, id ExecutionId) error {
	result, err := c.StopExecution(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StopExecution: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StopExecution: %+v", err)
	}

	return nil
}

// preparerForStopExecution prepares the StopExecution request.
func (c JobsClient) preparerForStopExecution(ctx context.Context, id ExecutionId) (*http.Request, error) {
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

// senderForStopExecution sends the StopExecution request. The method will close the
// http.Response Body if it receives an error.
func (c JobsClient) senderForStopExecution(ctx context.Context, req *http.Request) (future StopExecutionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
