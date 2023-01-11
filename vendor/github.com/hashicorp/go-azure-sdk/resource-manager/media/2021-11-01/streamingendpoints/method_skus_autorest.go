package streamingendpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusOperationResponse struct {
	HttpResponse *http.Response
	Model        *StreamingEndpointSkuInfoListResult
}

// Skus ...
func (c StreamingEndpointsClient) Skus(ctx context.Context, id StreamingEndpointId) (result SkusOperationResponse, err error) {
	req, err := c.preparerForSkus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "Skus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "Skus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSkus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingendpoints.StreamingEndpointsClient", "Skus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSkus prepares the Skus request.
func (c StreamingEndpointsClient) preparerForSkus(ctx context.Context, id StreamingEndpointId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSkus handles the response to the Skus request. The method always
// closes the http.Response Body.
func (c StreamingEndpointsClient) responderForSkus(resp *http.Response) (result SkusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
