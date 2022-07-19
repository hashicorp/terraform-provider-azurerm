package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideosUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoEntity
}

// VideosUpdate ...
func (c VideoAnalyzerClient) VideosUpdate(ctx context.Context, id VideoId, input VideoEntity) (result VideosUpdateOperationResponse, err error) {
	req, err := c.preparerForVideosUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideosUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideosUpdate prepares the VideosUpdate request.
func (c VideoAnalyzerClient) preparerForVideosUpdate(ctx context.Context, id VideoId, input VideoEntity) (*http.Request, error) {
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

// responderForVideosUpdate handles the response to the VideosUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosUpdate(resp *http.Response) (result VideosUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
