package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzersUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoAnalyzer
}

// VideoAnalyzersUpdate ...
func (c VideoAnalyzerClient) VideoAnalyzersUpdate(ctx context.Context, id VideoAnalyzerId, input VideoAnalyzerUpdate) (result VideoAnalyzersUpdateOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersUpdate prepares the VideoAnalyzersUpdate request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersUpdate(ctx context.Context, id VideoAnalyzerId, input VideoAnalyzerUpdate) (*http.Request, error) {
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

// responderForVideoAnalyzersUpdate handles the response to the VideoAnalyzersUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersUpdate(resp *http.Response) (result VideoAnalyzersUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
