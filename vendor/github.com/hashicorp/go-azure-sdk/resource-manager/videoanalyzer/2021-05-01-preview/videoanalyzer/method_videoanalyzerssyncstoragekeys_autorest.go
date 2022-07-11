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

type VideoAnalyzersSyncStorageKeysOperationResponse struct {
	HttpResponse *http.Response
}

// VideoAnalyzersSyncStorageKeys ...
func (c VideoAnalyzerClient) VideoAnalyzersSyncStorageKeys(ctx context.Context, id VideoAnalyzerId, input SyncStorageKeysInput) (result VideoAnalyzersSyncStorageKeysOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersSyncStorageKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersSyncStorageKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersSyncStorageKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersSyncStorageKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersSyncStorageKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersSyncStorageKeys prepares the VideoAnalyzersSyncStorageKeys request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersSyncStorageKeys(ctx context.Context, id VideoAnalyzerId, input SyncStorageKeysInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/syncStorageKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVideoAnalyzersSyncStorageKeys handles the response to the VideoAnalyzersSyncStorageKeys request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersSyncStorageKeys(resp *http.Response) (result VideoAnalyzersSyncStorageKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
