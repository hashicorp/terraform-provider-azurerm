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

type PausePrimingJobOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PausePrimingJob ...
func (c CachesClient) PausePrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (result PausePrimingJobOperationResponse, err error) {
	req, err := c.preparerForPausePrimingJob(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "PausePrimingJob", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPausePrimingJob(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "PausePrimingJob", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PausePrimingJobThenPoll performs PausePrimingJob then polls until it's completed
func (c CachesClient) PausePrimingJobThenPoll(ctx context.Context, id CacheId, input PrimingJobIdParameter) error {
	result, err := c.PausePrimingJob(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PausePrimingJob: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PausePrimingJob: %+v", err)
	}

	return nil
}

// preparerForPausePrimingJob prepares the PausePrimingJob request.
func (c CachesClient) preparerForPausePrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/pausePrimingJob", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPausePrimingJob sends the PausePrimingJob request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForPausePrimingJob(ctx context.Context, req *http.Request) (future PausePrimingJobOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
