package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoAnalyzer
}

// VideoAnalyzersGet ...
func (c VideoAnalyzerClient) VideoAnalyzersGet(ctx context.Context, id VideoAnalyzerId) (result VideoAnalyzersGetOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersGet prepares the VideoAnalyzersGet request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersGet(ctx context.Context, id VideoAnalyzerId) (*http.Request, error) {
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

// responderForVideoAnalyzersGet handles the response to the VideoAnalyzersGet request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersGet(resp *http.Response) (result VideoAnalyzersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
