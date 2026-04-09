package subscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingAccountGetPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *BillingAccountPoliciesResponse
}

// BillingAccountGetPolicy ...
func (c SubscriptionsClient) BillingAccountGetPolicy(ctx context.Context, id BillingAccountId) (result BillingAccountGetPolicyOperationResponse, err error) {
	req, err := c.preparerForBillingAccountGetPolicy(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "BillingAccountGetPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "BillingAccountGetPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForBillingAccountGetPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "BillingAccountGetPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForBillingAccountGetPolicy prepares the BillingAccountGetPolicy request.
func (c SubscriptionsClient) preparerForBillingAccountGetPolicy(ctx context.Context, id BillingAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Subscription/policies/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForBillingAccountGetPolicy handles the response to the BillingAccountGetPolicy request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForBillingAccountGetPolicy(resp *http.Response) (result BillingAccountGetPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
