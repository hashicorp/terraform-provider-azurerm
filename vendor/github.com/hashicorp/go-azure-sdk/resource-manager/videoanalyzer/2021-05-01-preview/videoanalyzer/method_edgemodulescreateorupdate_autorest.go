package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeModulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *EdgeModuleEntity
}

// EdgeModulesCreateOrUpdate ...
func (c VideoAnalyzerClient) EdgeModulesCreateOrUpdate(ctx context.Context, id EdgeModuleId, input EdgeModuleEntity) (result EdgeModulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForEdgeModulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEdgeModulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "EdgeModulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEdgeModulesCreateOrUpdate prepares the EdgeModulesCreateOrUpdate request.
func (c VideoAnalyzerClient) preparerForEdgeModulesCreateOrUpdate(ctx context.Context, id EdgeModuleId, input EdgeModuleEntity) (*http.Request, error) {
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

// responderForEdgeModulesCreateOrUpdate handles the response to the EdgeModulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForEdgeModulesCreateOrUpdate(resp *http.Response) (result EdgeModulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
