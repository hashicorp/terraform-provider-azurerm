package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPoliciesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AccessPoliciesDelete ...
func (c VideoAnalyzerClient) AccessPoliciesDelete(ctx context.Context, id AccessPoliciesId) (result AccessPoliciesDeleteOperationResponse, err error) {
	req, err := c.preparerForAccessPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccessPoliciesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccessPoliciesDelete prepares the AccessPoliciesDelete request.
func (c VideoAnalyzerClient) preparerForAccessPoliciesDelete(ctx context.Context, id AccessPoliciesId) (*http.Request, error) {
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

// responderForAccessPoliciesDelete handles the response to the AccessPoliciesDelete request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForAccessPoliciesDelete(resp *http.Response) (result AccessPoliciesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
