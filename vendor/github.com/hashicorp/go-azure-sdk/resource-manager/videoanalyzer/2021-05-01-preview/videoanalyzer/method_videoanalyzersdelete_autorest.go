package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzersDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// VideoAnalyzersDelete ...
func (c VideoAnalyzerClient) VideoAnalyzersDelete(ctx context.Context, id VideoAnalyzerId) (result VideoAnalyzersDeleteOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersDelete prepares the VideoAnalyzersDelete request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersDelete(ctx context.Context, id VideoAnalyzerId) (*http.Request, error) {
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

// responderForVideoAnalyzersDelete handles the response to the VideoAnalyzersDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersDelete(resp *http.Response) (result VideoAnalyzersDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
