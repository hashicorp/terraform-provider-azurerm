package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessPolicyEntity
}

// AccessPoliciesGet ...
func (c VideoAnalyzerClient) AccessPoliciesGet(ctx context.Context, id AccessPoliciesId) (result AccessPoliciesGetOperationResponse, err error) {
	req, err := c.preparerForAccessPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccessPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccessPoliciesGet prepares the AccessPoliciesGet request.
func (c VideoAnalyzerClient) preparerForAccessPoliciesGet(ctx context.Context, id AccessPoliciesId) (*http.Request, error) {
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

// responderForAccessPoliciesGet handles the response to the AccessPoliciesGet request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForAccessPoliciesGet(resp *http.Response) (result AccessPoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
