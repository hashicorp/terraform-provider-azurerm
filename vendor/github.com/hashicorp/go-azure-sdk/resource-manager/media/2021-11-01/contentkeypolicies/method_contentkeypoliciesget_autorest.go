package contentkeypolicies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContentKeyPolicy
}

// ContentKeyPoliciesGet ...
func (c ContentKeyPoliciesClient) ContentKeyPoliciesGet(ctx context.Context, id ContentKeyPolicyId) (result ContentKeyPoliciesGetOperationResponse, err error) {
	req, err := c.preparerForContentKeyPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContentKeyPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContentKeyPoliciesGet prepares the ContentKeyPoliciesGet request.
func (c ContentKeyPoliciesClient) preparerForContentKeyPoliciesGet(ctx context.Context, id ContentKeyPolicyId) (*http.Request, error) {
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

// responderForContentKeyPoliciesGet handles the response to the ContentKeyPoliciesGet request. The method always
// closes the http.Response Body.
func (c ContentKeyPoliciesClient) responderForContentKeyPoliciesGet(resp *http.Response) (result ContentKeyPoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
