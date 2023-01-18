package contentkeypolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesGetPolicyPropertiesWithSecretsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ContentKeyPolicyProperties
}

// ContentKeyPoliciesGetPolicyPropertiesWithSecrets ...
func (c ContentKeyPoliciesClient) ContentKeyPoliciesGetPolicyPropertiesWithSecrets(ctx context.Context, id ContentKeyPolicyId) (result ContentKeyPoliciesGetPolicyPropertiesWithSecretsOperationResponse, err error) {
	req, err := c.preparerForContentKeyPoliciesGetPolicyPropertiesWithSecrets(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGetPolicyPropertiesWithSecrets", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGetPolicyPropertiesWithSecrets", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForContentKeyPoliciesGetPolicyPropertiesWithSecrets(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesGetPolicyPropertiesWithSecrets", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForContentKeyPoliciesGetPolicyPropertiesWithSecrets prepares the ContentKeyPoliciesGetPolicyPropertiesWithSecrets request.
func (c ContentKeyPoliciesClient) preparerForContentKeyPoliciesGetPolicyPropertiesWithSecrets(ctx context.Context, id ContentKeyPolicyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getPolicyPropertiesWithSecrets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForContentKeyPoliciesGetPolicyPropertiesWithSecrets handles the response to the ContentKeyPoliciesGetPolicyPropertiesWithSecrets request. The method always
// closes the http.Response Body.
func (c ContentKeyPoliciesClient) responderForContentKeyPoliciesGetPolicyPropertiesWithSecrets(resp *http.Response) (result ContentKeyPoliciesGetPolicyPropertiesWithSecretsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
