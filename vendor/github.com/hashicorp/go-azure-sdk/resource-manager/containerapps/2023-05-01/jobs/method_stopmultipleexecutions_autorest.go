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

type StopMultipleExecutionsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StopMultipleExecutions ...
func (c JobsClient) StopMultipleExecutions(ctx context.Context, id JobId) (result StopMultipleExecutionsOperationResponse, err error) {
	req, err := c.preparerForStopMultipleExecutions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "StopMultipleExecutions", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStopMultipleExecutions(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "StopMultipleExecutions", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StopMultipleExecutionsThenPoll performs StopMultipleExecutions then polls until it's completed
func (c JobsClient) StopMultipleExecutionsThenPoll(ctx context.Context, id JobId) error {
	result, err := c.StopMultipleExecutions(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StopMultipleExecutions: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StopMultipleExecutions: %+v", err)
	}

	return nil
}

// preparerForStopMultipleExecutions prepares the StopMultipleExecutions request.
func (c JobsClient) preparerForStopMultipleExecutions(ctx context.Context, id JobId) (*http.Request, error) {
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

// senderForStopMultipleExecutions sends the StopMultipleExecutions request. The method will close the
// http.Response Body if it receives an error.
func (c JobsClient) senderForStopMultipleExecutions(ctx context.Context, req *http.Request) (future StopMultipleExecutionsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
