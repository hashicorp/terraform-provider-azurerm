package contentkeypolicies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContentKeyPolicy
}

// ContentKeyPoliciesUpdate ...
func (c ContentKeyPoliciesClient) ContentKeyPoliciesUpdate(ctx context.Context, id ContentKeyPolicyId, input ContentKeyPolicy) (result ContentKeyPoliciesUpdateOperationResponse, err error) {
	req, err := c.preparerForContentKeyPoliciesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContentKeyPoliciesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContentKeyPoliciesUpdate prepares the ContentKeyPoliciesUpdate request.
func (c ContentKeyPoliciesClient) preparerForContentKeyPoliciesUpdate(ctx context.Context, id ContentKeyPolicyId, input ContentKeyPolicy) (*http.Request, error) {
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

// responderForContentKeyPoliciesUpdate handles the response to the ContentKeyPoliciesUpdate request. The method always
// closes the http.Response Body.
func (c ContentKeyPoliciesClient) responderForContentKeyPoliciesUpdate(resp *http.Response) (result ContentKeyPoliciesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
