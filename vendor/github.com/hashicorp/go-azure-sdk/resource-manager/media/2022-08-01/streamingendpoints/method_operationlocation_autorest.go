package streamingendpoints

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
	Model        *StreamingEndpoint
}

// OperationLocation ...
func (c StreamingEndpointsClient) OperationLocation(ctx context.Context, id StreamingEndpointOperationLocationId) (result OperationLocationOperationResponse, err error) {
	req, err := c.preparerForOperationLocation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "OperationLocation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "OperationLocation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForOperationLocation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "OperationLocation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForOperationLocation prepares the OperationLocation request.
func (c StreamingEndpointsClient) preparerForOperationLocation(ctx context.Context, id StreamingEndpointOperationLocationId) (*http.Request, error) {
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
func (c StreamingEndpointsClient) responderForOperationLocation(resp *http.Response) (result OperationLocationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
