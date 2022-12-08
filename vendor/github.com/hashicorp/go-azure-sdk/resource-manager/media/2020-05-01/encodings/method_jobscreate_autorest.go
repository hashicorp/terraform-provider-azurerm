package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Job
}

// JobsCreate ...
func (c EncodingsClient) JobsCreate(ctx context.Context, id JobId, input Job) (result JobsCreateOperationResponse, err error) {
	req, err := c.preparerForJobsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobsCreate prepares the JobsCreate request.
func (c EncodingsClient) preparerForJobsCreate(ctx context.Context, id JobId, input Job) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobsCreate handles the response to the JobsCreate request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForJobsCreate(resp *http.Response) (result JobsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
