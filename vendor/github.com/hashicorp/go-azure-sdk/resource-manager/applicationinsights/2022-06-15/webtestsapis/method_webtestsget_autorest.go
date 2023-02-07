package webtestsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *WebTest
}

// WebTestsGet ...
func (c WebTestsAPIsClient) WebTestsGet(ctx context.Context, id WebTestId) (result WebTestsGetOperationResponse, err error) {
	req, err := c.preparerForWebTestsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWebTestsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWebTestsGet prepares the WebTestsGet request.
func (c WebTestsAPIsClient) preparerForWebTestsGet(ctx context.Context, id WebTestId) (*http.Request, error) {
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

// responderForWebTestsGet handles the response to the WebTestsGet request. The method always
// closes the http.Response Body.
func (c WebTestsAPIsClient) responderForWebTestsGet(resp *http.Response) (result WebTestsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
