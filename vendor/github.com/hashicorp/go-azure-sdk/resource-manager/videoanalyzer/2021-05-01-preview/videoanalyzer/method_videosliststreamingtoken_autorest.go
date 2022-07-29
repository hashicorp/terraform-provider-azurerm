package videoanalyzer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideosListStreamingTokenOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoStreamingToken
}

// VideosListStreamingToken ...
func (c VideoAnalyzerClient) VideosListStreamingToken(ctx context.Context, id VideoId) (result VideosListStreamingTokenOperationResponse, err error) {
	req, err := c.preparerForVideosListStreamingToken(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosListStreamingToken", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosListStreamingToken", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideosListStreamingToken(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosListStreamingToken", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideosListStreamingToken prepares the VideosListStreamingToken request.
func (c VideoAnalyzerClient) preparerForVideosListStreamingToken(ctx context.Context, id VideoId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listStreamingToken", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVideosListStreamingToken handles the response to the VideosListStreamingToken request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosListStreamingToken(resp *http.Response) (result VideosListStreamingTokenOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
