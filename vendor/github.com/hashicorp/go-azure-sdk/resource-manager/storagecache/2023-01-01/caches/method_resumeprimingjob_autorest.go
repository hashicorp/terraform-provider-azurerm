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

type ResumePrimingJobOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResumePrimingJob ...
func (c CachesClient) ResumePrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (result ResumePrimingJobOperationResponse, err error) {
	req, err := c.preparerForResumePrimingJob(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "ResumePrimingJob", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResumePrimingJob(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "caches.CachesClient", "ResumePrimingJob", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResumePrimingJobThenPoll performs ResumePrimingJob then polls until it's completed
func (c CachesClient) ResumePrimingJobThenPoll(ctx context.Context, id CacheId, input PrimingJobIdParameter) error {
	result, err := c.ResumePrimingJob(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ResumePrimingJob: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResumePrimingJob: %+v", err)
	}

	return nil
}

// preparerForResumePrimingJob prepares the ResumePrimingJob request.
func (c CachesClient) preparerForResumePrimingJob(ctx context.Context, id CacheId, input PrimingJobIdParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resumePrimingJob", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResumePrimingJob sends the ResumePrimingJob request. The method will close the
// http.Response Body if it receives an error.
func (c CachesClient) senderForResumePrimingJob(ctx context.Context, req *http.Request) (future ResumePrimingJobOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
