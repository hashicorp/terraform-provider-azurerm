package webapplicationfirewallpolicies

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *WebApplicationFirewallPolicy
}

// PoliciesGet ...
func (c WebApplicationFirewallPoliciesClient) PoliciesGet(ctx context.Context, id FrontDoorWebApplicationFirewallPolicyId) (result PoliciesGetOperationResponse, err error) {
	req, err := c.preparerForPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPoliciesGet prepares the PoliciesGet request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesGet(ctx context.Context, id FrontDoorWebApplicationFirewallPolicyId) (*http.Request, error) {
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

// responderForPoliciesGet handles the response to the PoliciesGet request. The method always
// closes the http.Response Body.
func (c WebApplicationFirewallPoliciesClient) responderForPoliciesGet(resp *http.Response) (result PoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
