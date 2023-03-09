package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// JobsDelete ...
func (c EncodingsClient) JobsDelete(ctx context.Context, id JobId) (result JobsDeleteOperationResponse, err error) {
	req, err := c.preparerForJobsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobsDelete prepares the JobsDelete request.
func (c EncodingsClient) preparerForJobsDelete(ctx context.Context, id JobId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobsDelete handles the response to the JobsDelete request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForJobsDelete(resp *http.Response) (result JobsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
