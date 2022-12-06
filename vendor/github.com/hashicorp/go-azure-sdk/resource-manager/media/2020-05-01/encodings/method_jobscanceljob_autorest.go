package encodings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsCancelJobOperationResponse struct {
	HttpResponse *http.Response
}

// JobsCancelJob ...
func (c EncodingsClient) JobsCancelJob(ctx context.Context, id JobId) (result JobsCancelJobOperationResponse, err error) {
	req, err := c.preparerForJobsCancelJob(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCancelJob", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCancelJob", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForJobsCancelJob(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsCancelJob", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForJobsCancelJob prepares the JobsCancelJob request.
func (c EncodingsClient) preparerForJobsCancelJob(ctx context.Context, id JobId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cancelJob", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForJobsCancelJob handles the response to the JobsCancelJob request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForJobsCancelJob(resp *http.Response) (result JobsCancelJobOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
