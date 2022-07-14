package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideosGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoEntity
}

// VideosGet ...
func (c VideoAnalyzerClient) VideosGet(ctx context.Context, id VideoId) (result VideosGetOperationResponse, err error) {
	req, err := c.preparerForVideosGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideosGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideosGet prepares the VideosGet request.
func (c VideoAnalyzerClient) preparerForVideosGet(ctx context.Context, id VideoId) (*http.Request, error) {
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

// responderForVideosGet handles the response to the VideosGet request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosGet(resp *http.Response) (result VideosGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
