package liveoutputs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AsyncOperationOperationResponse struct {
	HttpResponse *http.Response
	Model        *AsyncOperationResult
}

// AsyncOperation ...
func (c LiveOutputsClient) AsyncOperation(ctx context.Context, id LiveOutputOperationId) (result AsyncOperationOperationResponse, err error) {
	req, err := c.preparerForAsyncOperation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "AsyncOperation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "AsyncOperation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAsyncOperation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "AsyncOperation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAsyncOperation prepares the AsyncOperation request.
func (c LiveOutputsClient) preparerForAsyncOperation(ctx context.Context, id LiveOutputOperationId) (*http.Request, error) {
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

// responderForAsyncOperation handles the response to the AsyncOperation request. The method always
// closes the http.Response Body.
func (c LiveOutputsClient) responderForAsyncOperation(resp *http.Response) (result AsyncOperationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
