package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideosCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoEntity
}

// VideosCreateOrUpdate ...
func (c VideoAnalyzerClient) VideosCreateOrUpdate(ctx context.Context, id VideoId, input VideoEntity) (result VideosCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForVideosCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideosCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideosCreateOrUpdate prepares the VideosCreateOrUpdate request.
func (c VideoAnalyzerClient) preparerForVideosCreateOrUpdate(ctx context.Context, id VideoId, input VideoEntity) (*http.Request, error) {
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

// responderForVideosCreateOrUpdate handles the response to the VideosCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosCreateOrUpdate(resp *http.Response) (result VideosCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
