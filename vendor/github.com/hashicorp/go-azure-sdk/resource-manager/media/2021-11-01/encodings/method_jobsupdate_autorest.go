package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Job
}

// JobsUpdate ...
func (c EncodingsClient) JobsUpdate(ctx context.Context, id JobId, input Job) (result JobsUpdateOperationResponse, err error) {
	req, err := c.preparerForJobsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobsUpdate prepares the JobsUpdate request.
func (c EncodingsClient) preparerForJobsUpdate(ctx context.Context, id JobId, input Job) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobsUpdate handles the response to the JobsUpdate request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForJobsUpdate(resp *http.Response) (result JobsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
