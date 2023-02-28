package endpoints

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsEndpointGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DigitalTwinsEndpointResource
}

// DigitalTwinsEndpointGet ...
func (c EndpointsClient) DigitalTwinsEndpointGet(ctx context.Context, id EndpointId) (result DigitalTwinsEndpointGetOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsEndpointGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDigitalTwinsEndpointGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDigitalTwinsEndpointGet prepares the DigitalTwinsEndpointGet request.
func (c EndpointsClient) preparerForDigitalTwinsEndpointGet(ctx context.Context, id EndpointId) (*http.Request, error) {
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

// responderForDigitalTwinsEndpointGet handles the response to the DigitalTwinsEndpointGet request. The method always
// closes the http.Response Body.
func (c EndpointsClient) responderForDigitalTwinsEndpointGet(resp *http.Response) (result DigitalTwinsEndpointGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
