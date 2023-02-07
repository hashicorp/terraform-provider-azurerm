package liveoutputs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationLocationOperationResponse struct {
	HttpResponse *http.Response
	Model        *LiveOutput
}

// OperationLocation ...
func (c LiveOutputsClient) OperationLocation(ctx context.Context, id LiveOutputOperationLocationId) (result OperationLocationOperationResponse, err error) {
	req, err := c.preparerForOperationLocation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "OperationLocation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "OperationLocation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForOperationLocation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveoutputs.LiveOutputsClient", "OperationLocation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForOperationLocation prepares the OperationLocation request.
func (c LiveOutputsClient) preparerForOperationLocation(ctx context.Context, id LiveOutputOperationLocationId) (*http.Request, error) {
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

// responderForOperationLocation handles the response to the OperationLocation request. The method always
// closes the http.Response Body.
func (c LiveOutputsClient) responderForOperationLocation(resp *http.Response) (result OperationLocationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
