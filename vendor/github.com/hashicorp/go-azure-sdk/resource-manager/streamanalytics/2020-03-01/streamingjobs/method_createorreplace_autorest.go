package streamingjobs

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

type CreateOrReplaceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type CreateOrReplaceOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOrReplaceOperationOptions() CreateOrReplaceOperationOptions {
	return CreateOrReplaceOperationOptions{}
}

func (o CreateOrReplaceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	if o.IfNoneMatch != nil {
		out["If-None-Match"] = *o.IfNoneMatch
	}

	return out
}

func (o CreateOrReplaceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// CreateOrReplace ...
func (c StreamingJobsClient) CreateOrReplace(ctx context.Context, id StreamingJobId, input StreamingJob, options CreateOrReplaceOperationOptions) (result CreateOrReplaceOperationResponse, err error) {
	req, err := c.preparerForCreateOrReplace(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingjobs.StreamingJobsClient", "CreateOrReplace", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateOrReplace(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingjobs.StreamingJobsClient", "CreateOrReplace", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateOrReplaceThenPoll performs CreateOrReplace then polls until it's completed
func (c StreamingJobsClient) CreateOrReplaceThenPoll(ctx context.Context, id StreamingJobId, input StreamingJob, options CreateOrReplaceOperationOptions) error {
	result, err := c.CreateOrReplace(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing CreateOrReplace: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateOrReplace: %+v", err)
	}

	return nil
}

// preparerForCreateOrReplace prepares the CreateOrReplace request.
func (c StreamingJobsClient) preparerForCreateOrReplace(ctx context.Context, id StreamingJobId, input StreamingJob, options CreateOrReplaceOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateOrReplace sends the CreateOrReplace request. The method will close the
// http.Response Body if it receives an error.
func (c StreamingJobsClient) senderForCreateOrReplace(ctx context.Context, req *http.Request) (future CreateOrReplaceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
