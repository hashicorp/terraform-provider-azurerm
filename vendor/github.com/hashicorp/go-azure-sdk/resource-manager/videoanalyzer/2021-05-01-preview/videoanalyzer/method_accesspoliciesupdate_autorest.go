package videoanalyzer

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPoliciesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessPolicyEntity
}

// AccessPoliciesUpdate ...
func (c VideoAnalyzerClient) AccessPoliciesUpdate(ctx context.Context, id AccessPoliciesId, input AccessPolicyEntity) (result AccessPoliciesUpdateOperationResponse, err error) {
	req, err := c.preparerForAccessPoliciesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccessPoliciesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccessPoliciesUpdate prepares the AccessPoliciesUpdate request.
func (c VideoAnalyzerClient) preparerForAccessPoliciesUpdate(ctx context.Context, id AccessPoliciesId, input AccessPolicyEntity) (*http.Request, error) {
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

// responderForAccessPoliciesUpdate handles the response to the AccessPoliciesUpdate request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForAccessPoliciesUpdate(resp *http.Response) (result AccessPoliciesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
