package updateruns

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

type StartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
	Model        *UpdateRun
}

type StartOperationOptions struct {
	IfMatch *string
}

func DefaultStartOperationOptions() StartOperationOptions {
	return StartOperationOptions{}
}

func (o StartOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o StartOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Start ...
func (c UpdateRunsClient) Start(ctx context.Context, id UpdateRunId, options StartOperationOptions) (result StartOperationResponse, err error) {
	req, err := c.preparerForStart(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "Start", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "Start", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StartThenPoll performs Start then polls until it's completed
func (c UpdateRunsClient) StartThenPoll(ctx context.Context, id UpdateRunId, options StartOperationOptions) error {
	result, err := c.Start(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Start: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Start: %+v", err)
	}

	return nil
}

// preparerForStart prepares the Start request.
func (c UpdateRunsClient) preparerForStart(ctx context.Context, id UpdateRunId, options StartOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/start", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStart sends the Start request. The method will close the
// http.Response Body if it receives an error.
func (c UpdateRunsClient) senderForStart(ctx context.Context, req *http.Request) (future StartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
