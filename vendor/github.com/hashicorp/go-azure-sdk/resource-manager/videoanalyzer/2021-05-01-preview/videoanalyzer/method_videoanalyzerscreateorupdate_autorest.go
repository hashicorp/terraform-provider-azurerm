package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzersCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoAnalyzer
}

// VideoAnalyzersCreateOrUpdate ...
func (c VideoAnalyzerClient) VideoAnalyzersCreateOrUpdate(ctx context.Context, id VideoAnalyzerId, input VideoAnalyzer) (result VideoAnalyzersCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersCreateOrUpdate prepares the VideoAnalyzersCreateOrUpdate request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersCreateOrUpdate(ctx context.Context, id VideoAnalyzerId, input VideoAnalyzer) (*http.Request, error) {
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

// responderForVideoAnalyzersCreateOrUpdate handles the response to the VideoAnalyzersCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersCreateOrUpdate(resp *http.Response) (result VideoAnalyzersCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
