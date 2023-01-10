package webtestsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsUpdateTagsOperationResponse struct {
	HttpResponse *http.Response
	Model        *WebTest
}

// WebTestsUpdateTags ...
func (c WebTestsAPIsClient) WebTestsUpdateTags(ctx context.Context, id WebTestId, input TagsResource) (result WebTestsUpdateTagsOperationResponse, err error) {
	req, err := c.preparerForWebTestsUpdateTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsUpdateTags", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsUpdateTags", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWebTestsUpdateTags(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsUpdateTags", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWebTestsUpdateTags prepares the WebTestsUpdateTags request.
func (c WebTestsAPIsClient) preparerForWebTestsUpdateTags(ctx context.Context, id WebTestId, input TagsResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWebTestsUpdateTags handles the response to the WebTestsUpdateTags request. The method always
// closes the http.Response Body.
func (c WebTestsAPIsClient) responderForWebTestsUpdateTags(resp *http.Response) (result WebTestsUpdateTagsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
