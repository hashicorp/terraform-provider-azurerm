package globalschedules

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

type ExecuteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Execute ...
func (c GlobalSchedulesClient) Execute(ctx context.Context, id ScheduleId) (result ExecuteOperationResponse, err error) {
	req, err := c.preparerForExecute(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "globalschedules.GlobalSchedulesClient", "Execute", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExecute(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "globalschedules.GlobalSchedulesClient", "Execute", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExecuteThenPoll performs Execute then polls until it's completed
func (c GlobalSchedulesClient) ExecuteThenPoll(ctx context.Context, id ScheduleId) error {
	result, err := c.Execute(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Execute: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Execute: %+v", err)
	}

	return nil
}

// preparerForExecute prepares the Execute request.
func (c GlobalSchedulesClient) preparerForExecute(ctx context.Context, id ScheduleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/execute", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExecute sends the Execute request. The method will close the
// http.Response Body if it receives an error.
func (c GlobalSchedulesClient) senderForExecute(ctx context.Context, req *http.Request) (future ExecuteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
