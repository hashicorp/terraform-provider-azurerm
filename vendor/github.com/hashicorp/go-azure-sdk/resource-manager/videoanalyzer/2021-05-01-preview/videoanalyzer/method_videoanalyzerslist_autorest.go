package videoanalyzer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzersListOperationResponse struct {
	HttpResponse *http.Response
	Model        *VideoAnalyzerCollection
}

// VideoAnalyzersList ...
func (c VideoAnalyzerClient) VideoAnalyzersList(ctx context.Context, id commonids.ResourceGroupId) (result VideoAnalyzersListOperationResponse, err error) {
	req, err := c.preparerForVideoAnalyzersList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVideoAnalyzersList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideoAnalyzersList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVideoAnalyzersList prepares the VideoAnalyzersList request.
func (c VideoAnalyzerClient) preparerForVideoAnalyzersList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Media/videoAnalyzers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVideoAnalyzersList handles the response to the VideoAnalyzersList request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideoAnalyzersList(resp *http.Response) (result VideoAnalyzersListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
