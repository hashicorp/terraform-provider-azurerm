package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionPolicyGetPolicyForTenantOperationResponse struct {
	HttpResponse *http.Response
	Model        *GetTenantPolicyResponse
}

// SubscriptionPolicyGetPolicyForTenant ...
func (c SubscriptionsClient) SubscriptionPolicyGetPolicyForTenant(ctx context.Context) (result SubscriptionPolicyGetPolicyForTenantOperationResponse, err error) {
	req, err := c.preparerForSubscriptionPolicyGetPolicyForTenant(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyGetPolicyForTenant", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyGetPolicyForTenant", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionPolicyGetPolicyForTenant(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyGetPolicyForTenant", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionPolicyGetPolicyForTenant prepares the SubscriptionPolicyGetPolicyForTenant request.
func (c SubscriptionsClient) preparerForSubscriptionPolicyGetPolicyForTenant(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Subscription/policies/default"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionPolicyGetPolicyForTenant handles the response to the SubscriptionPolicyGetPolicyForTenant request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionPolicyGetPolicyForTenant(resp *http.Response) (result SubscriptionPolicyGetPolicyForTenantOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
