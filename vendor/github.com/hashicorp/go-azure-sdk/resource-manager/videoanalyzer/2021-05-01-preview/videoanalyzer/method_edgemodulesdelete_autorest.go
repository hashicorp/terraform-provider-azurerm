package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// EdgeModulesDelete ...
func (c VideoAnalyzerClient) EdgeModulesDelete(ctx context.Context, id EdgeModuleId) (result EdgeModulesDeleteOperationResponse, err error) {
	req, err := c.preparerForEdgeModulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEdgeModulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEdgeModulesDelete prepares the EdgeModulesDelete request.
func (c VideoAnalyzerClient) preparerForEdgeModulesDelete(ctx context.Context, id EdgeModuleId) (*http.Request, error) {
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

// responderForEdgeModulesDelete handles the response to the EdgeModulesDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForEdgeModulesDelete(resp *http.Response) (result EdgeModulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
