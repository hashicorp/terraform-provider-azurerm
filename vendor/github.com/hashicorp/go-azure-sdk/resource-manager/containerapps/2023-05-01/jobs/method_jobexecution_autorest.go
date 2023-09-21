package jobs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionOperationResponse struct {
	HttpResponse *http.Response
	Model        *JobExecution
}

// JobExecution ...
func (c JobsClient) JobExecution(ctx context.Context, id ExecutionId) (result JobExecutionOperationResponse, err error) {
	req, err := c.preparerForJobExecution(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "JobExecution", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "JobExecution", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobExecution(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "JobExecution", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobExecution prepares the JobExecution request.
func (c JobsClient) preparerForJobExecution(ctx context.Context, id ExecutionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobExecution handles the response to the JobExecution request. The method always
// closes the http.Response Body.
func (c JobsClient) responderForJobExecution(resp *http.Response) (result JobExecutionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
