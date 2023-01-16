package webtestsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// WebTestsDelete ...
func (c WebTestsAPIsClient) WebTestsDelete(ctx context.Context, id WebTestId) (result WebTestsDeleteOperationResponse, err error) {
	req, err := c.preparerForWebTestsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWebTestsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWebTestsDelete prepares the WebTestsDelete request.
func (c WebTestsAPIsClient) preparerForWebTestsDelete(ctx context.Context, id WebTestId) (*http.Request, error) {
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

// responderForWebTestsDelete handles the response to the WebTestsDelete request. The method always
// closes the http.Response Body.
func (c WebTestsAPIsClient) responderForWebTestsDelete(resp *http.Response) (result WebTestsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
