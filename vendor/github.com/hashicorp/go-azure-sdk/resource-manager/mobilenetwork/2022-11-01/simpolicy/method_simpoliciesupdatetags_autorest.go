package simpolicy

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimPoliciesUpdateTagsOperationResponse struct {
	HttpResponse *http.Response
	Model        *SimPolicy
}

// SimPoliciesUpdateTags ...
func (c SIMPolicyClient) SimPoliciesUpdateTags(ctx context.Context, id SimPolicyId, input TagsObject) (result SimPoliciesUpdateTagsOperationResponse, err error) {
	req, err := c.preparerForSimPoliciesUpdateTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesUpdateTags", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesUpdateTags", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSimPoliciesUpdateTags(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesUpdateTags", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSimPoliciesUpdateTags prepares the SimPoliciesUpdateTags request.
func (c SIMPolicyClient) preparerForSimPoliciesUpdateTags(ctx context.Context, id SimPolicyId, input TagsObject) (*http.Request, error) {
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

// responderForSimPoliciesUpdateTags handles the response to the SimPoliciesUpdateTags request. The method always
// closes the http.Response Body.
func (c SIMPolicyClient) responderForSimPoliciesUpdateTags(resp *http.Response) (result SimPoliciesUpdateTagsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
