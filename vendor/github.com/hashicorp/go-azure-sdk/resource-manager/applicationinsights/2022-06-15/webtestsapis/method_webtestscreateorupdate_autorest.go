package webtestsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *WebTest
}

// WebTestsCreateOrUpdate ...
func (c WebTestsAPIsClient) WebTestsCreateOrUpdate(ctx context.Context, id WebTestId, input WebTest) (result WebTestsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForWebTestsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWebTestsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWebTestsCreateOrUpdate prepares the WebTestsCreateOrUpdate request.
func (c WebTestsAPIsClient) preparerForWebTestsCreateOrUpdate(ctx context.Context, id WebTestId, input WebTest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWebTestsCreateOrUpdate handles the response to the WebTestsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c WebTestsAPIsClient) responderForWebTestsCreateOrUpdate(resp *http.Response) (result WebTestsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
