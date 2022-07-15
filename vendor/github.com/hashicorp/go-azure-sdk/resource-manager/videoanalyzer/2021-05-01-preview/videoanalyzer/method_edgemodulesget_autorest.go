package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *EdgeModuleEntity
}

// EdgeModulesGet ...
func (c VideoAnalyzerClient) EdgeModulesGet(ctx context.Context, id EdgeModuleId) (result EdgeModulesGetOperationResponse, err error) {
	req, err := c.preparerForEdgeModulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEdgeModulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEdgeModulesGet prepares the EdgeModulesGet request.
func (c VideoAnalyzerClient) preparerForEdgeModulesGet(ctx context.Context, id EdgeModuleId) (*http.Request, error) {
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

// responderForEdgeModulesGet handles the response to the EdgeModulesGet request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForEdgeModulesGet(resp *http.Response) (result EdgeModulesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
