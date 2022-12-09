package simpolicy

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimPoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SimPolicy
}

// SimPoliciesGet ...
func (c SIMPolicyClient) SimPoliciesGet(ctx context.Context, id SimPolicyId) (result SimPoliciesGetOperationResponse, err error) {
	req, err := c.preparerForSimPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSimPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "simpolicy.SIMPolicyClient", "SimPoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSimPoliciesGet prepares the SimPoliciesGet request.
func (c SIMPolicyClient) preparerForSimPoliciesGet(ctx context.Context, id SimPolicyId) (*http.Request, error) {
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

// responderForSimPoliciesGet handles the response to the SimPoliciesGet request. The method always
// closes the http.Response Body.
func (c SIMPolicyClient) responderForSimPoliciesGet(resp *http.Response) (result SimPoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
