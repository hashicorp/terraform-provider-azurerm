package frontdoors

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendEndpointsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *FrontendEndpoint
}

// FrontendEndpointsGet ...
func (c FrontDoorsClient) FrontendEndpointsGet(ctx context.Context, id FrontendEndpointId) (result FrontendEndpointsGetOperationResponse, err error) {
	req, err := c.preparerForFrontendEndpointsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFrontendEndpointsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFrontendEndpointsGet prepares the FrontendEndpointsGet request.
func (c FrontDoorsClient) preparerForFrontendEndpointsGet(ctx context.Context, id FrontendEndpointId) (*http.Request, error) {
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

// responderForFrontendEndpointsGet handles the response to the FrontendEndpointsGet request. The method always
// closes the http.Response Body.
func (c FrontDoorsClient) responderForFrontendEndpointsGet(resp *http.Response) (result FrontendEndpointsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
