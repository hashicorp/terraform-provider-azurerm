package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionPolicyAddUpdatePolicyForTenantOperationResponse struct {
	HttpResponse *http.Response
	Model        *GetTenantPolicyResponse
}

// SubscriptionPolicyAddUpdatePolicyForTenant ...
func (c SubscriptionsClient) SubscriptionPolicyAddUpdatePolicyForTenant(ctx context.Context, input PutTenantPolicyRequestProperties) (result SubscriptionPolicyAddUpdatePolicyForTenantOperationResponse, err error) {
	req, err := c.preparerForSubscriptionPolicyAddUpdatePolicyForTenant(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyAddUpdatePolicyForTenant", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyAddUpdatePolicyForTenant", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionPolicyAddUpdatePolicyForTenant(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionPolicyAddUpdatePolicyForTenant", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionPolicyAddUpdatePolicyForTenant prepares the SubscriptionPolicyAddUpdatePolicyForTenant request.
func (c SubscriptionsClient) preparerForSubscriptionPolicyAddUpdatePolicyForTenant(ctx context.Context, input PutTenantPolicyRequestProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Subscription/policies/default"),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionPolicyAddUpdatePolicyForTenant handles the response to the SubscriptionPolicyAddUpdatePolicyForTenant request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionPolicyAddUpdatePolicyForTenant(resp *http.Response) (result SubscriptionPolicyAddUpdatePolicyForTenantOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
