package caches

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

type StopPrimingJobOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StopPrimingJob ...
func (c CachesClient) StopPrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (result StopPrimingJobOperationResponse, err error) {
	req, err := c.preparerForStopPrimingJob(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "StopPrimingJob", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStopPrimingJob(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "StopPrimingJob", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StopPrimingJobThenPoll performs StopPrimingJob then polls until it's completed
func (c CachesClient) StopPrimingJobThenPoll(ctx context.Context, id CacheId, input PrimingJobIdParameter) error {
	result, err := c.StopPrimingJob(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing StopPrimingJob: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StopPrimingJob: %+v", err)
	}

	return nil
}

// preparerForStopPrimingJob prepares the StopPrimingJob request.
func (c CachesClient) preparerForStopPrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stopPrimingJob", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStopPrimingJob sends the StopPrimingJob request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForStopPrimingJob(ctx context.Context, req *http.Request) (future StopPrimingJobOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
